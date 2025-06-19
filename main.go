package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	//db *database.Queries
	DB_URL    string
	PLATFORM  string
	JWT_TOKEN string
	PORT      string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	PLATFORM := os.Getenv("PLATFORM")
	if PLATFORM == "" {
		log.Fatal("PLATFORM must be set")
	}

	JWT_TOKEN := os.Getenv("JWT_TOKEN")

	_ = apiConfig{
		DB_URL:    dbURL,
		PLATFORM:  PLATFORM,
		JWT_TOKEN: JWT_TOKEN,
		PORT:      port,
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
