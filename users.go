package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gcheong/Chirpy/internal/auth"
	"github.com/gcheong/Chirpy/internal/database"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)

	e := parameters{}
	err := decoder.Decode(&e)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	hash, err := auth.HashPassword(e.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          e.Email,
		HashedPassword: hash,
	})

	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)

	e := parameters{}
	err := decoder.Decode(&e)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), e.Email)
	if err != nil {
		log.Printf("Error getting user by email: %s", err)
		w.WriteHeader(500)
		return
	}

	valid, err := auth.CheckPasswordHash(e.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error verifying password: %s", err)
		w.WriteHeader(401)
		return
	}
	if !valid {
		log.Printf("Invalid password for user: %s", e.Email)
		w.WriteHeader(401)
		return
	}

	// token, err := auth.GenerateJWT(user.ID.String())
	// if err != nil {
	// 	log.Printf("Error generating JWT: %s", err)
	// 	w.WriteHeader(500)
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		// Token: token,
	})
}
