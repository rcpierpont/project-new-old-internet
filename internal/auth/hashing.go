package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPW, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		fmt.Printf("unable to hash password: %s\n", err)
		return "", err
	}
	return hashedPW, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Printf("unable to compare password with hash: %s\n", err)
	}
	return match, nil
}
