package handlers

import (
	"encoding/json"
	"net/http"
)

// SendResponse returns our status along with our proper non-corsy headers
func SendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

	jsonData, err := json.Marshal(data)
	if err != nil {
		SendResponse(w, err.Error(), 500)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
