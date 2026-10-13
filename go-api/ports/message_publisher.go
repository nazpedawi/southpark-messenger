package ports

import "southpark/go-api/domain"

// MessagePublisher is the OUTBOUND (secondary) port.
// The core uses it to send a message out of the application.
// It does not know or care that RabbitMQ is behind it.
type MessagePublisher interface {
	Publish(message domain.Message) *domain.AppError
}