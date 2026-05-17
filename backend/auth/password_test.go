package auth

import (
	"testing"
)

func TestHashAndCheckPassword_RoundTrip(t *testing.T) {
	password := "my-secure-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty hash")
	}

	if !CheckPassword(hash, password) {
		t.Error("CheckPassword should return true for correct password, but returned false")
	}
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if CheckPassword(hash, wrongPassword) {
		t.Error("CheckPassword should return false for wrong password, but returned true")
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "same-password"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword (first call) returned error: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword (second call) returned error: %v", err)
	}

	if hash1 == hash2 {
		t.Error("HashPassword should produce different hashes for the same password due to bcrypt salt, but got identical hashes")
	}

	// Both hashes should still validate against the original password.
	if !CheckPassword(hash1, password) {
		t.Error("CheckPassword failed for first hash")
	}
	if !CheckPassword(hash2, password) {
		t.Error("CheckPassword failed for second hash")
	}
}
