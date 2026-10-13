package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"southpark/go-api/adapters/rabbitmq"
	"southpark/go-api/adapters/rest"
	"southpark/go-api/app"

	"github.com/gorilla/mux"
)

const webPort = ":8080"

func main() {
	fmt.Println("Starting App")

	var router = mux.NewRouter()

	// Where is RabbitMQ? Locally: localhost. In Docker Compose we set RABBITMQ_URL.
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	// wiring: publisher (outbound adapter) -> service -> handlers (inbound adapter) -> routes
	publisher, err := rabbitmq.NewRabbitMQPublisher(rabbitURL)
	if err != nil {
		log.Fatal("Could not connect to RabbitMQ: ", err)
	}
	defer publisher.Close()

	messageService := app.NewMessageService(publisher)
	messageHandlers := rest.MessageHandlers{Service: messageService}

	router.HandleFunc("/messages", messageHandlers.SendMessage).
		Methods("POST").
		Name("SendMessage")

	fmt.Println("Starting Web Server on port", webPort)
	log.Fatal(http.ListenAndServe(webPort, router))
}