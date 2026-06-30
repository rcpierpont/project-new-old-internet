package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {
	authz := strings.Replace(headers.Get("Authorization"), "Bearer ", "", 1)
	if len(authz) == 0 {
		return "", fmt.Errorf("not authorized")
	}
	return authz, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	key := []byte(tokenSecret)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(key)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("error validating JWT: %v\n", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing subject from claims: %v\n", err)
	}

	subjectUUID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing UUID from subject: %v\n", err)
	}

	return subjectUUID, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	key := strings.Replace(authHeader, "ApiKey ", "", 1)
	if len(key) == 0 {
		return "", fmt.Errorf("no api key found in auth header: %s\n", authHeader)
	}
	return key, nil
}
