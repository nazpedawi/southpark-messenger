package ports

import "southpark/go-api/domain"

// MessageService is the INBOUND (primary) port.
// It describes what the application can do.
// The REST adapter calls these methods; it never sees the real implementation.
type MessageService interface {
	SendMessage(author string, body string) (*domain.Message, *domain.AppError)
}