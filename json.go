package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErrors(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(&payload)

	if err != nil {
		log.Println("Failed to encode json response")
		w.WriteHeader(500)
		return
	}
}
