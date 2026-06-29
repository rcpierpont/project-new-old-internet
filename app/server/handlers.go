package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *envConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode user parameters", err)
		return
	}

	fmt.Printf("creating user: %s\n", params.Email)
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "unable to create user", err)
		return
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		response{User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		},
	)
}

func (cfg *envConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}
	userIDPathVal := r.PathValue("userID")

	log.Printf("user id: %s", userIDPathVal)
	userID, err := uuid.Parse(userIDPathVal)
	if err != nil {
		userIDPathVal = ""
	}

	if userIDPathVal == "" {
		log.Fatal("User not found")
		return
	}

	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to get user from provided id", err)
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		response{User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		},
	)
}
