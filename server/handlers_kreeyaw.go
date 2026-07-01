package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rcpierpont/project-new-old-internet/internal/auth"
	"github.com/rcpierpont/project-new-old-internet/internal/database"
)

type errorResponse struct {
	Error string `json:"error"`
}

type kreeyawParams struct {
	Body string `json:"body"`
}

type kreeyawResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

const maxChars int = 255

func (cfg *envConfig) handlerCreateKreeyaw(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil || userID == uuid.Nil {
		respondWithError(w, http.StatusUnauthorized, "unable to authenticate user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := kreeyawParams{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters", err)
		return
	}

	if len(params.Body) > maxChars {
		log.Printf("Error, body too long\n")
		respondWithError(w, http.StatusBadRequest, "length must be 140 chars or less", err)
		return
	}

	chirp, err := cfg.db.CreateKreeyaw(r.Context(), database.CreateKreeyawParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error, unable to create chirp: %s\n", err)
	}
	respondWithJSON(w, http.StatusCreated, kreeyawResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
