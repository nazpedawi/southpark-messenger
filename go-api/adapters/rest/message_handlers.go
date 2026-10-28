package rest

import (
	"encoding/json"
	"net/http"
	"southpark/go-api/domain"
	"southpark/go-api/ports"
)

// We define the handlers by wrapping them in a struct.
// MessageHandlers contains the message service (the inbound port), and we create the receiver functions (handlers) on it.
type MessageHandlers struct {
	Service ports.MessageService
}

// sendMessageRequest is the shape of the JSON the client sends
type sendMessageRequest struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

// SendMessage handles POST /messages
func (mh *MessageHandlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	// 1. Turn the JSON body into a Go struct
	var request sendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		appErr := domain.ValidationError("invalid JSON: " + err.Error())
		writeResponse(w, appErr.Code, appErr)
		return
	}

	// 2. Call the core through the inbound port
	message, appErr := mh.Service.SendMessage(request.Author, request.Body)
	if appErr != nil {
		writeResponse(w, appErr.Code, appErr)
		return
	}

	// 3. Success: 201 Created, and send the message back
	writeResponse(w, http.StatusCreated, message)
}
