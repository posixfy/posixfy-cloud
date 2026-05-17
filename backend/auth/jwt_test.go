package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key-for-unit-tests"

func TestGenerateToken_CreatesValidToken(t *testing.T) {
	tokenStr, err := GenerateToken(testSecret, 42, "alice", 1000, 1000, "1000,1001", "admin")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("GenerateToken returned empty token string")
	}
}

func TestParseToken_ValidToken(t *testing.T) {
	tokenStr, err := GenerateToken(testSecret, 42, "alice", 1000, 1000, "1000,1001", "admin")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	claims, err := ParseToken(testSecret, tokenStr)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("expected UserID 42, got %d", claims.UserID)
	}
	if claims.Username != "alice" {
		t.Errorf("expected Username 'alice', got %q", claims.Username)
	}
	if claims.UID != 1000 {
		t.Errorf("expected UID 1000, got %d", claims.UID)
	}
	if claims.GID != 1000 {
		t.Errorf("expected GID 1000, got %d", claims.GID)
	}
	if claims.Groups != "1000,1001" {
		t.Errorf("expected Groups '1000,1001', got %q", claims.Groups)
	}
	if claims.Role != "admin" {
		t.Errorf("expected Role 'admin', got %q", claims.Role)
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// Manually craft a token that expired in the past.
	claims := Claims{
		UserID:   1,
		Username: "bob",
		UID:      500,
		GID:      500,
		Groups:   "500",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	_, err = ParseToken(testSecret, tokenStr)
	if err == nil {
		t.Fatal("ParseToken should reject expired token, but returned nil error")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	tokenStr, err := GenerateToken(testSecret, 1, "carol", 1000, 1000, "1000", "user")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	_, err = ParseToken("wrong-secret", tokenStr)
	if err == nil {
		t.Fatal("ParseToken should reject token signed with different secret, but returned nil error")
	}
}

func TestParseToken_MalformedToken(t *testing.T) {
	malformedTokens := []string{
		"",
		"not-a-jwt",
		"three.parts.here",
		"eyJhbGciOiJIUzI1NiJ9.invalid.payload",
	}

	for _, tok := range malformedTokens {
		_, err := ParseToken(testSecret, tok)
		if err == nil {
			t.Errorf("ParseToken should reject malformed token %q, but returned nil error", tok)
		}
	}
}
