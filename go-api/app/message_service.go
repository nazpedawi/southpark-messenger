package app

import (
	"southpark/go-api/domain"
	"southpark/go-api/ports"
)

// DefaultMessageService is the real implementation of the ports.MessageService interface (the inbound port).
// We pass in the publisher (an interface) as a dependency.
type DefaultMessageService struct {
	publisher ports.MessagePublisher
}

// SendMessage is a receiver function: it attaches the method to DefaultMessageService, which makes it satisfy ports.MessageService.
func (s DefaultMessageService) SendMessage(author string, body string) (*domain.Message, *domain.AppError) {
	// 1. Let the domain create the message (and check the rules)
	message, err := domain.NewMessage(author, body)
	if err != nil {
		return nil, err
	}

	// 2. Send it out through the outbound port.
	// we talk to the interface, not to RabbitMQ directly.
	publishErr := s.publisher.Publish(*message)
	if publishErr != nil {
		return nil, publishErr
	}

	return message, nil
}

// NewMessageService is a helper function to instantiate the message service.
func NewMessageService(publisher ports.MessagePublisher) DefaultMessageService {
	return DefaultMessageService{publisher}
}

// Compile-time check: the build fails if DefaultMessageService ever stops satisfying the ports.MessageService interface.
var _ ports.MessageService = DefaultMessageService{}