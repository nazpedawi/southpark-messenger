package domain

import (
	"strings"
	"time"
)

// Message is the core entity
// The json tags control the key names when it is turned into JSON.
type Message struct {
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// NewMessage is a "constructor" helper: the proper way to create a Message.
// It enforces the business rules, so an invalid Message can never exist.
func NewMessage(author, body string) (*Message, *AppError) {
	author = strings.TrimSpace(author) 
	body = strings.TrimSpace(body)

	if author == "" {
		return nil, ValidationError("author is required")
	}
	if body == "" {
		return nil, ValidationError("body is required")
	}

	return &Message{
		Author:    author,
		Body:      body,
		CreatedAt: time.Now(),
	}, nil // nil means: no error
}