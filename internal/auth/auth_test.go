package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	validPassword := "X9f$2pL!zW"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate valid hash: %v", err)
	}

	cases := []struct {
		name      string
		inputPass string
		exitHash  string
		wantErr   bool
	}{
		{
			name:      "Valid Password",
			inputPass: validPassword,
			exitHash:  validHash,
			wantErr:   false,
		},
		{
			name:      "Invalid Hash",
			inputPass: validPassword,
			exitHash:  "valid$Password",
			wantErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Logf("Testing hash: %q", c.exitHash)
			gotErr := VerifyPassword(c.exitHash, c.inputPass)
			if (gotErr != nil) != c.wantErr {
				t.Errorf("Test Failed => expected: %v is not the same as actual: %v", c.wantErr, err)
				return
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validString := "supersecretkey123456"
	invalidString := "totallywrongsecret"
	token, err := MakeJWT(userID, validString, time.Minute*15)
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err)
	}
	cases := []struct {
		name        string
		inputString string
		inputToken  string
		wantErr     bool
	}{
		{
			name:        "Valid String",
			inputString: validString,
			inputToken:  token,
			wantErr:     false,
		},
		{
			name:        "Invalid String",
			inputString: invalidString,
			inputToken:  token,
			wantErr:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotID, err := ValidateJWT(c.inputToken, c.inputString)
			if (err != nil) != c.wantErr {
				t.Errorf("Test Failed => expected: %v, got: %v", c.wantErr, err)
				return
			}
			if err == nil && gotID != userID {
				t.Errorf("Test Failed => expected: %v, got: %v", userID, gotID)
			}
		})
	}
}
