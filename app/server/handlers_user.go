package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rcpierpont/project-new-old-internet/internal/auth"
	"github.com/rcpierpont/project-new-old-internet/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}

func (cfg *envConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	})
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
