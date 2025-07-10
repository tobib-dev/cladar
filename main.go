package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tobib-dev/cladar/internal/database"
)

type apiConfig struct {
	db        *database.Queries
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

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	PLATFORM := os.Getenv("PLATFORM")
	if PLATFORM == "" {
		log.Fatal("PLATFORM must be set")
	}

	JWT_TOKEN := os.Getenv("JWT_TOKEN")
	if JWT_TOKEN == "" {
		log.Fatalf("JWT TOKEN must be set")
	}

	cfg := apiConfig{
		db:        dbQueries,
		DB_URL:    dbURL,
		PLATFORM:  PLATFORM,
		JWT_TOKEN: JWT_TOKEN,
		PORT:      port,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	printBanner()
	log.Printf("Serving on port: %s\n", port)

	mux.HandleFunc("POST /api/customers", cfg.handlerCreateCustomer)
	mux.HandleFunc("GET /api/customers", cfg.handlerGetAllCustomers)
	mux.HandleFunc("GET /api/customers/{custID}", cfg.handlerGetCustomer)
	mux.HandleFunc("PUT /api/customers/{custID}", cfg.handlerUpdateCustomer)
	mux.HandleFunc("DELETE /api/customers/{custID}", cfg.handlerDeleteCustomer)

	mux.HandleFunc("POST /api/agents", cfg.handlerCreateAgent)
	mux.HandleFunc("GET /api/agents", cfg.handlerGetAllAgents)
	mux.HandleFunc("GET /api/agents/{agentID}", cfg.handlerGetAgent)
	mux.HandleFunc("PUT /api/agents/{agentID}", cfg.handlerUpdateAgents)
	mux.HandleFunc("DELETE /api/agents/{agentID}", cfg.handlerDeleteAgent)

	mux.HandleFunc("POST /api/claims/{custID}", cfg.handlerCreateClaim)
	mux.HandleFunc("GET /api/claims", cfg.handlerGetAllClaims)
	mux.HandleFunc("GET /api/claims/{claimID}", cfg.handlerGetClaim)
	mux.HandleFunc("GET /api/customers/{custID}/claims", cfg.handlerGetClaimsByCustomer)
	mux.HandleFunc("GET /api/agents/{agentID}/claims", cfg.handlerGetClaimsByAssignedAgent)
	mux.HandleFunc("PUT /api/claims/change_agents", cfg.handlerChangeAssignedAgent)
	mux.HandleFunc("PUT /api/claims/update_claimtype/{claimID}", cfg.handlerChangeClaimType)
	mux.HandleFunc("PUT /api/claims/decline_claims/{claimID}", cfg.handlerDeclineClaim)
	mux.HandleFunc("PUT /api/awards/award_claims/{claimID}", cfg.handlerAwardClaim)
	mux.HandleFunc("PUT /api/awards/update_amount_for_awards/{claimID}", cfg.handlerChangeAwardAmount)
	mux.HandleFunc("PUT /api/awards/complete_claims/{claimID}", cfg.handlerCompleteClaim)
	mux.HandleFunc("DELETE /api/claims/{claimID}", cfg.handlerDeleteClaim)

	// Get claims by status
	mux.HandleFunc("GET /api/claims/pending", cfg.handlerGetPendingClaims)
	mux.HandleFunc("GET /api/claims/awarded", cfg.handlerGetAwardedClaims)
	mux.HandleFunc("GET /api/claims/declined", cfg.handlerGetDeclinedClaims)

	mux.HandleFunc("POST /api/departments", cfg.handlerCreateDept)
	mux.HandleFunc("GET /api/departments", cfg.handlerGetAllDepts)

	mux.HandleFunc("POST /api/managers", cfg.handlerCreateManager)
	mux.HandleFunc("GET /api/managers", cfg.handlerGetAllManagers)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshTokens)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevokeTokens)

	//mux.HandleFunc("POST /api/reset", cfg.handlerReset)
	log.Fatal(srv.ListenAndServe())
}

func printBanner() {
	fmt.Print(`
   ________          __
  / ____/ /___ _____/ /___ ______
 / /   / / __  / __  / __  / ___/
/ /___/ / /_/ / /_/ / /_/ / /
\____/_/\__,_/\__,_/\__,_/_/
`)
}
