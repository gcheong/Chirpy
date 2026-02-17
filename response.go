package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {	
	dat, err := json.Marshal(payload)
	if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
	}
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    w.Write(dat)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, chirpError{Err: message})
}