package auth

import (
	"net/http"
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
				t.Errorf("Test Failed => expected: %v is not the same as actual: %v", true, c.wantErr)
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
				t.Errorf("Test Failed => expected: %v, got: %v", true, c.wantErr)
				return
			}
			if err == nil && gotID != userID {
				t.Errorf("Test Failed => expected: %v, got: %v", userID, gotID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	validHeader := http.Header{}
	validHeader.Set("Authorization", "Bearer supersecretkey123456")
	validString, err := GetBearerToken(validHeader)
	if err != nil {
		t.Fatalf("Failed to get bearer token")
		return
	}

	noBearer := http.Header{}
	noBearer.Set("Authorization", "supersecretkey123456")

	noSpace := http.Header{}
	noSpace.Set("Authorization", "Bearersupersecretkey123456")

	cases := []struct {
		name        string
		inputHeader http.Header
		tokenString string
		wantErr     bool
	}{
		{
			name:        "Valid Header",
			inputHeader: validHeader,
			tokenString: validString,
			wantErr:     false,
		},
		{
			name:        "No Bearer tag",
			inputHeader: noBearer,
			tokenString: validString,
			wantErr:     true,
		},
		{
			name:        "No space between bearer & tag",
			inputHeader: noSpace,
			tokenString: validString,
			wantErr:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := GetBearerToken(c.inputHeader)
			if (err != nil) != c.wantErr {
				t.Errorf("Test Failed, expected: %v, got: %v", true, c.wantErr)
				return
			}
		})
	}
}
