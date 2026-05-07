// Package config loads environment variables for the TutorAI backend.
package config

import (
	"os"
	"time"
)

// Config holds all runtime configuration for the TutorAI backend.
type Config struct {
	Port              string
	OllamaBaseURL     string
	OllamaLLMModel    string
	DataServiceURL    string
	DataServiceAPIKey string
	// HTTPTimeout is applied to the shared http.Client used for all outbound
	// calls (Ollama and data service). Local LLM inference can be slow, so
	// this is set generously.
	HTTPTimeout time.Duration
}

// Load reads configuration from environment variables.
// All variables are expected to be set before calling Load (e.g. via godotenv).
func Load() Config {
	return Config{
		Port:              getEnv("PORT", "8000"),
		OllamaBaseURL:     getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		OllamaLLMModel:    getEnv("OLLAMA_LLM_MODEL", "llama3.1"),
		DataServiceURL:    getEnv("DATA_SERVICE_URL", "http://localhost:8001"),
		DataServiceAPIKey: getEnv("DATA_SERVICE_API_KEY", ""),
		HTTPTimeout:       120 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
