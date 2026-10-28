package app

import (
	"southpark/go-api/domain"
	"testing"
)

// fakePublisher pretends to be RabbitMQ. Because it has a Publish method with the right signature, it satisfies ports.MessagePublisher.
type fakePublisher struct {
	published []domain.Message // remembers what was "sent"
	failWith  *domain.AppError // if set, Publish returns this error
}

func (f *fakePublisher) Publish(message domain.Message) *domain.AppError {
	if f.failWith != nil {
		return f.failWith
	}
	f.published = append(f.published, message)
	return nil
}

func TestSendMessage_PublishesValidMessage(t *testing.T) {
	fake := &fakePublisher{}
	service := NewMessageService(fake)

	message, err := service.SendMessage("Cartman", "Respect my authoritah!")

	if err != nil {
		t.Fatalf("expected no error, got %q", err.AsMessage())
	}
	if message.Author != "Cartman" {
		t.Errorf("expected author Cartman, got %q", message.Author)
	}
	if len(fake.published) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(fake.published))
	}
}

func TestSendMessage_RejectsMissingAuthor(t *testing.T) {
	fake := &fakePublisher{}
	service := NewMessageService(fake)

	_, err := service.SendMessage("   ", "Oh my God, they killed Kenny!")

	if err == nil || err.Code != 400 {
		t.Fatalf("expected a 400 error for a missing author")
	}
	if len(fake.published) != 0 {
		t.Errorf("an invalid message must never be published")
	}
}

func TestSendMessage_RejectsMissingBody(t *testing.T) {
	fake := &fakePublisher{}
	service := NewMessageService(fake)

	_, err := service.SendMessage("Kyle", "")

	if err == nil || err.Code != 400 {
		t.Fatalf("expected a 400 error for a missing body")
	}
}

func TestSendMessage_ReturnsPublisherError(t *testing.T) {
	fake := &fakePublisher{failWith: domain.UnexpectedError("rabbitmq is down")}
	service := NewMessageService(fake)

	_, err := service.SendMessage("Butters", "Oh hamburgers.")

	if err == nil || err.Code != 500 {
		t.Fatalf("expected a 500 error when publishing fails")
	}
}
