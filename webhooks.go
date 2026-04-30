package main

import (
	"encoding/json"
	"log"
	"net/http"	

	"github.com/gcheong/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type webhook_data struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data webhook_data `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)

	e := parameters{}
	err := decoder.Decode(&e)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}	

	if e.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	api_key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid api key", err)
		return
	}

	if api_key != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Mismatched api key", err)
		return
	}

	_, err = cfg.dbQueries.UpgradeUser(r.Context(), e.Data.UserID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not upgrade user", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}