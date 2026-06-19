package config

import (
	"os"
	"strings"
)

type Config struct {
	FilebridgeURL     string
	APIKey            string
	JWTSecret         string
	DBPath            string
	ListenAddr        string
	AdminInitPassword string
	CORSOrigins       []string
	LogLevel          string
}

func Load() *Config {
	return &Config{
		FilebridgeURL:     getEnv("FILEBRIDGE_URL", "http://127.0.0.1:3000"),
		APIKey:            mustEnv("API_KEY"),
		JWTSecret:         mustEnv("JWT_SECRET"),
		DBPath:            getEnv("DB_PATH", "./posixfy-cloud.db"),
		ListenAddr:        getEnv("LISTEN_ADDR", "0.0.0.0:8080"),
		AdminInitPassword: os.Getenv("ADMIN_INIT_PASSWORD"),
		CORSOrigins:       parseCORSOrigins(getEnv("CORS_ORIGINS", "*")),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}
}

func parseCORSOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	var origins []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	return origins
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required env var " + key + " is not set")
	}
	return v
}
