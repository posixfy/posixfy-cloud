package config

import (
	"os"
	"reflect"
	"testing"
)

// setRequiredEnvVars sets the mandatory env vars that Load() requires via mustEnv.
// It returns a cleanup function that restores the original values.
func setRequiredEnvVars(t *testing.T) func() {
	t.Helper()

	origAPIKey := os.Getenv("API_KEY")
	origJWTSecret := os.Getenv("JWT_SECRET")

	os.Setenv("API_KEY", "test-api-key")
	os.Setenv("JWT_SECRET", "test-jwt-secret")

	return func() {
		if origAPIKey == "" {
			os.Unsetenv("API_KEY")
		} else {
			os.Setenv("API_KEY", origAPIKey)
		}
		if origJWTSecret == "" {
			os.Unsetenv("JWT_SECRET")
		} else {
			os.Setenv("JWT_SECRET", origJWTSecret)
		}
	}
}

func TestParseCORSOrigins_Wildcard(t *testing.T) {
	cleanup := setRequiredEnvVars(t)
	defer cleanup()

	origCORS := os.Getenv("CORS_ORIGINS")
	os.Setenv("CORS_ORIGINS", "*")
	defer func() {
		if origCORS == "" {
			os.Unsetenv("CORS_ORIGINS")
		} else {
			os.Setenv("CORS_ORIGINS", origCORS)
		}
	}()

	cfg := Load()

	expected := []string{"*"}
	if !reflect.DeepEqual(cfg.CORSOrigins, expected) {
		t.Errorf("expected CORSOrigins %v, got %v", expected, cfg.CORSOrigins)
	}
}

func TestParseCORSOrigins_MultipleOrigins(t *testing.T) {
	cleanup := setRequiredEnvVars(t)
	defer cleanup()

	origCORS := os.Getenv("CORS_ORIGINS")
	os.Setenv("CORS_ORIGINS", "http://localhost:5173,https://example.com,https://app.example.com")
	defer func() {
		if origCORS == "" {
			os.Unsetenv("CORS_ORIGINS")
		} else {
			os.Setenv("CORS_ORIGINS", origCORS)
		}
	}()

	cfg := Load()

	expected := []string{"http://localhost:5173", "https://example.com", "https://app.example.com"}
	if !reflect.DeepEqual(cfg.CORSOrigins, expected) {
		t.Errorf("expected CORSOrigins %v, got %v", expected, cfg.CORSOrigins)
	}
}

func TestParseCORSOrigins_WithWhitespace(t *testing.T) {
	cleanup := setRequiredEnvVars(t)
	defer cleanup()

	origCORS := os.Getenv("CORS_ORIGINS")
	os.Setenv("CORS_ORIGINS", " http://localhost:5173 , https://example.com , https://app.example.com ")
	defer func() {
		if origCORS == "" {
			os.Unsetenv("CORS_ORIGINS")
		} else {
			os.Setenv("CORS_ORIGINS", origCORS)
		}
	}()

	cfg := Load()

	expected := []string{"http://localhost:5173", "https://example.com", "https://app.example.com"}
	if !reflect.DeepEqual(cfg.CORSOrigins, expected) {
		t.Errorf("expected CORSOrigins %v (whitespace trimmed), got %v", expected, cfg.CORSOrigins)
	}
}
