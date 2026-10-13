package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"southpark/go-api/domain"
	"southpark/go-api/ports"
)

// The names of our RabbitMQ "plumbing".
const (
	ExchangeName = "southpark"          // topic exchange: the "post office"
	QueueName    = "southpark_messages" // the "mailbox"
	BindingKey   = "messages.#"         // rule: every routing key starting with "messages."
)

// We wrap the RabbitMQ stuff in a struct RabbitMQPublisher.
// It is the OUTBOUND adapter: it implements the ports.MessagePublisher interface.
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex // a channel must not be used by two requests at the same time
}

// NewRabbitMQPublisher is the "constructor" helper.
// It connects to RabbitMQ and sets up the exchange, the queue and the binding.
func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	// Step 1: dial the connection (with retries: RabbitMQ in Docker needs a few seconds to start)
	conn, err := dialWithRetry(url, 10, 3*time.Second)
	if err != nil {
		return nil, err
	}

	// Step 2: open the channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Step 3a: declare the TOPIC exchange
	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type: route on routing-key patterns
		true,         // durable (survives a RabbitMQ restart)
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Step 3b: declare the queue
	_, err = ch.QueueDeclare(
		QueueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Step 3c: bind the queue to the exchange with the binding key
	err = ch.QueueBind(
		QueueName,    // queue name
		BindingKey,   // binding key "messages.#"
		ExchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("RabbitMQ: connected, exchange/queue/binding ready")
	return &RabbitMQPublisher{conn: conn, channel: ch}, nil
}

// Publish is a receiver function: it makes RabbitMQPublisher satisfy ports.MessagePublisher.
func (p *RabbitMQPublisher) Publish(message domain.Message) *domain.AppError {
	// Turn the domain message into JSON bytes
	body, err := json.Marshal(message)
	if err != nil {
		return domain.UnexpectedError("could not encode message")
	}

	// e.g. "messages.cartman" -> matches the binding "messages.#"
	routingKey := "messages." + routingKeyPart(message.Author)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 4: publish the message
	p.mu.Lock()
	defer p.mu.Unlock()
	err = p.channel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Println("RabbitMQ: error publishing message: " + err.Error())
		return domain.UnexpectedError("could not publish message")
	}

	log.Printf("RabbitMQ: published to %q with routing key %q", ExchangeName, routingKey)
	return nil
}

// Close closes the channel and the connection.
func (p *RabbitMQPublisher) Close() {
	p.channel.Close()
	p.conn.Close()
}

// dialWithRetry tries to connect a few times before giving up.
func dialWithRetry(url string, attempts int, wait time.Duration) (*amqp.Connection, error) {
	var err error
	for i := 1; i <= attempts; i++ {
		var conn *amqp.Connection
		conn, err = amqp.Dial(url)
		if err == nil {
			return conn, nil
		}
		log.Printf("RabbitMQ: connection attempt %d/%d failed, retrying in %s...", i, attempts, wait)
		time.Sleep(wait)
	}
	return nil, err
}

// routingKeyPart turns "Mr. Garrison" into "mr__garrison".
// Dots must go, because RabbitMQ uses dots to separate routing-key words.
func routingKeyPart(author string) string {
	lower := strings.ToLower(strings.TrimSpace(author))
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, lower)
}

// Compile-time check: RabbitMQPublisher must satisfy the outbound port.
var _ ports.MessagePublisher = (*RabbitMQPublisher)(nil)