package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	cases := []struct {
		name     string
		password string
		check    string
		want     bool
	}{
		{"matching", "thisisatest", "thisisatest", true},
		{"different", "thisisatest", "wrongpass", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hash, _ := HashPassword(c.password)
			match, err := CheckPasswordHash(c.check, hash)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if match != c.want {
				t.Errorf("got %v, want %v", match, c.want)
			}
		})
	}
}
