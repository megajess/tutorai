package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"tutorai/backend/config"
	"tutorai/backend/internal/api"
	"tutorai/backend/internal/retrieval"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env: %v", err)
	}

	cfg := config.Load()

	colorLookup, err := retrieval.LoadColorLookup("data/color_identity_lookup.json")
	if err != nil {
		log.Fatalf("failed to load color identity lookup: %v", err)
	}

	httpClient := &http.Client{Timeout: cfg.HTTPTimeout}
	dataClient := retrieval.NewClient(cfg, httpClient)
	chatHandler := api.NewChatHandler(cfg, httpClient, colorLookup, dataClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	r.Get("/health", handleHealth)
	r.Post("/chat", chatHandler.ServeHTTP)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
