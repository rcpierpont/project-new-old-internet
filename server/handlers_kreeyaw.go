package main

import (
	"encoding/json"
	"errors"
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

// POST must contain text body for kreeyaw
// requires valid bearer token because kreeyaw must be linked to author user ID in db
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

	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unexpected error parsing user from validated JWT", err)
	}
	log.Printf("authenticated user: %s\n", user.Email)

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

	kreeyaw, err := cfg.db.CreateKreeyaw(r.Context(), database.CreateKreeyawParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error, unable to create chirp: %s\n", err)
	}
	respondWithJSON(w, http.StatusCreated, kreeyawResponse{
		ID:        kreeyaw.ID,
		CreatedAt: kreeyaw.CreatedAt,
		UpdatedAt: kreeyaw.UpdatedAt,
		Body:      kreeyaw.Body,
		UserID:    kreeyaw.UserID,
	})
}

// targeted delete by providing kreeyaw id param
// requires valid bearer token owned by author of kreeyaw being deleted
func (cfg *envConfig) handlerDeleteKreeyaw(w http.ResponseWriter, r *http.Request) {
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

	kreeyawPathVal := r.PathValue("kreeyawID")
	kreeyawUUID, err := uuid.Parse(kreeyawPathVal)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "provided kreeyaw ID not found in database", err)
		return
	}

	kreeyaw, err := cfg.db.GetKreeyaw(r.Context(), kreeyawUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve kreeyaw using id from url param", err)
		return
	}

	if kreeyaw.UserID != userID {
		respondWithError(w, http.StatusForbidden, "unable to delete chirps owned by other users", errors.New("unauthorized delete request"))
		return
	}

	err = cfg.db.DeleteKreeyaw(r.Context(), kreeyawUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unexpected error trying to delete kreeyaw", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// retrieve single kreeyaw by providing kreeyaw id as url param
// must be authenticated user (must include valid bearer token)
func (cfg *envConfig) handlerGetKreeyawByID(w http.ResponseWriter, r *http.Request) {
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

	kreeyawPathVal := r.PathValue("kreeyawID")
	kreeyawUUID, err := uuid.Parse(kreeyawPathVal)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "provided kreeyaw ID not found in database", err)
		return
	}

	kreeyaw, err := cfg.db.GetKreeyaw(r.Context(), kreeyawUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve kreeyaw using id from url param", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, kreeyawResponse{
		ID:        kreeyaw.ID,
		CreatedAt: kreeyaw.CreatedAt,
		UpdatedAt: kreeyaw.UpdatedAt,
		Body:      kreeyaw.Body,
		UserID:    kreeyaw.UserID,
	})
}
