package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// Arrange
	userID := uuid.New()
	secret := "my-test-secret"

	// Act: create a token
	tokenString, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Act: validate it
	gotUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}

	// Assert
	if gotUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, gotUserID)
	}
}

func TestJWTExpiredToken(t *testing.T) {
	// Arrange
	userID := uuid.New()
	secret := "my-test-secret"

	// Act: use an expired token
	tokenString, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
}

func TestJWTWrongSecret(t *testing.T) {
	// Arrange
	userID := uuid.New()

	// Act: use an expired token
	tokenString, err := MakeJWT(userID, "secret-A", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Act: Validate it
	_, err = ValidateJWT(tokenString, "secret-B")
	if err == nil {
		t.Errorf("expected error for wrong secret, got nil")
	}
}

func TestGetBearerToken(t *testing.T) {
	// Arrange
	headers := http.Header{}
	expectedBearer := "abc123"
	headers.Set("Authorization", "Bearer " + expectedBearer)

	// Act: call GetBearerToken
	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned error: %v", err)
	}

	// Assert: returned header == expectedHeader
	if token != expectedBearer {
		t.Errorf("Expected bearer %v, got bearer %v", expectedBearer, token)
	}
}

func TestGetBearerToken_Missing(t *testing.T) {
	// Arrange
	headers := http.Header{}

	// Act: call GetBearerToken
	token, err := GetBearerToken(headers)
	if err == nil {
		t.Fatalf("expected an error for missing Auth header, got nil (token=%q)", token)
	}

	// Assert: returned header == expectedHeader
	if token != "" {
		t.Errorf("expected empty token on error, got %q", token)
	}
}
