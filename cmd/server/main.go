package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"e-commerce/internal/api"
	"e-commerce/internal/mailer"
	"e-commerce/internal/store"
)

func main() {
	st, err := store.Open(getenv("DB_PATH", "ecommerce.db"))
	if err != nil {
		log.Fatalf("ouverture de la connexion à la base de données : %v", err)
	}

	server := &http.Server{
		Addr:         getenv("LISTEN_ADDR", ":8080"),
		Handler:      api.NewRouter(st, mailer.LogMailer{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Serveur démarré sur %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
