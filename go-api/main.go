package main

import (
	"fmt"
	"log"
	"net/http"
	"southpark/go-api/adapters/rest"
	"southpark/go-api/app"
	"southpark/go-api/domain"
	"github.com/gorilla/mux"
)

const webPort = ":8080"

// TEMPORARY: a stand-in publisher that just logs the message.
// In Step 7 we replace it with the real RabbitMQ adapter.
type logPublisher struct{}

func (p logPublisher) Publish(message domain.Message) *domain.AppError {
	log.Printf("[temporary publisher] %s says: %s", message.Author, message.Body)
	return nil
}

func main() {
	fmt.Println("Starting App")

	var router = mux.NewRouter()

	// wiring: publisher -> service -> handlers -> routes
	publisher := logPublisher{}
	messageService := app.NewMessageService(publisher)
	messageHandlers := rest.MessageHandlers{Service: messageService}

	router.HandleFunc("/messages", messageHandlers.SendMessage).
		Methods("POST").
		Name("SendMessage")

	fmt.Println("Starting Web Server on port", webPort)
	log.Fatal(http.ListenAndServe(webPort, router))
}