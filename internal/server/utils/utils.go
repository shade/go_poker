package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorMsg struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, status int) {
	dataStr, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(status)
	w.Write(dataStr)
}
