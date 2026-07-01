package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rcpierpont/project-new-old-internet/internal/auth"
	"github.com/rcpierpont/project-new-old-internet/internal/database"
)

const refreshTokenExpiryDays int = 60
const accessTokenExpiryMinutes int = 60

func (cfg *envConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode login credentials", err)
		return
	}

	expiry := min(params.ExpiresInSeconds, 3600)
	if expiry == 0 {
		expiry = 3600
	}
	expiryDuration, err := time.ParseDuration(fmt.Sprintf("%ds", expiry))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid expiration duration, check format", err)
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to find user in database", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to verify password", err)
		return
	}
	if !match {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, expiryDuration)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating JWT", err)
	}

	refTokenData := auth.MakeRefreshToken()
	refExpiryDuration, err := time.ParseDuration(fmt.Sprintf("%dd", refreshTokenExpiryDays))
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refTokenData,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refExpiryDuration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token", err)
	}

	respondWithJSON(
		w,
		http.StatusOK, response{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
			Token:        token,
			RefreshToken: refreshToken.Token,
		},
	)
}

func (cfg *envConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", err)
		return
	}

	foundToken, err := cfg.db.CheckToken(r.Context(), refreshToken)
	userFromToken := foundToken.UserID
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized - unable to validate token", err)
		return
	}

	if foundToken.RevokedAt.Valid && foundToken.RevokedAt.Time.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "not authorized - token expired", err)
		return
	}

	tokenExpiry, err := time.ParseDuration(fmt.Sprintf("%dm", accessTokenExpiryMinutes))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unknown error while parsing duration - possible bad request", err)
		return
	}
	accessToken, err := auth.MakeJWT(userFromToken, cfg.secret, tokenExpiry)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to create jwt token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *envConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "no token found to revoke", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token: refreshToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "token not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
