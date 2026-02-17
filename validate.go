package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"strings"

)

func handlerValidate(w http.ResponseWriter, r *http.Request){
    type chirp struct {
  		Body string `json:"body"`
	}	

    decoder := json.NewDecoder(r.Body)
    c := chirp{}
    err := decoder.Decode(&c)
	log.Printf("Received chirp: %s", c.Body)
	fmt.Printf("Received chirp: %s", c.Body)
    
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
    }

	if len(c.Body) > 140 {
		log.Printf("Chirp is too long: %d characters", len(c.Body))
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	respondWithJSON(w, 200, chirpValid{CleanedBody: cleaneadText(c.Body)})
}

func cleaneadText(s string) string {
	words := strings.Split(s, " ")	
	var cleanWords []string
	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			cleanWords = append(cleanWords, "****")
		default:
			cleanWords = append(cleanWords, word)
			
		}
	}

	return strings.Join(cleanWords, " ")
}