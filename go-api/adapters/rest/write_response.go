package rest

import (
	"encoding/json"
	"net/http"
)

// writeResponse sets the JSON header, the status code and encodes the data.
// Every handler uses it, so the handlers stay short.
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}
}
