package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	validPassword := "X9f$2pL!zW"
	validHash, err := HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate valid hash: %v", err)
	}

	t.Logf("Generated hash: %q", validHash)
	t.Logf("Hash length: %d", len(validHash))

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
