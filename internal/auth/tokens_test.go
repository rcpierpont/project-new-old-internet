package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	testUUID := uuid.New()
	cases := []struct {
		name   string
		uuid   uuid.UUID
		secret string
		check  uuid.UUID
		want   uuid.UUID
	}{
		{"successful-validation", testUUID, "test-secret", testUUID, testUUID},
		{"unsuccessful-bad-secret", testUUID, "test-invalid-secret", testUUID, uuid.Nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			duration, err := time.ParseDuration("5m")
			if err != nil {
				t.Fatalf("unexpected error generating duration: %v", err)
			}
			tokenString, err := MakeJWT(c.uuid, "test-secret", duration)
			if err != nil {
				t.Fatalf("unexpected error creating JWT: %v", err)
			}
			uuid, err := ValidateJWT(tokenString, c.secret)
			if c.want != uuid && err != nil {
				t.Fatalf("unexpected error validating JWT: %v", err)
			}
			if uuid != c.want {
				t.Errorf("got %v, want %v", c.check, c.want)
			}
		})
	}
}
