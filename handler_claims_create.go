package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

type Claims struct {
	ID              uuid.UUID `json:"id"`
	CustomerID      uuid.UUID `json:"customer_id"`
	AssignedAgentID uuid.UUID `json:"assigned_agent_id"`
	ClaimType       string    `json:"claim_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CurrentStatus   string    `json:"current_status"`
	AwardAmount     string    `json:"award_amount"`
}

func (cfg *apiConfig) handlerCreateClaim(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		AgentIDString string `json:"agent_id"`
		ClaimType     string `json:"claim_type"`
	}

	type Response struct {
		Claims
	}

	custIDString := r.PathValue("custID")
	custID, err := uuid.Parse(custIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get customerID", err)
		return
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	agentID, err := uuid.Parse(params.AgentIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get agentID as uuid", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't retrieve bearer token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.UserRole != database.UserType(UserRoleManager) {
		respondWithError(w, http.StatusUnauthorized, "Only managers can create claims", nil)
		return
	}

	claim, err := cfg.db.CreateClaim(r.Context(), database.CreateClaimParams{
		CustomerID: custID,
		AgentID:    agentID,
		ClaimType:  params.ClaimType,
		// Assign null value to award during creation. Award will be updated when awarded
		Award: sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create award", err)
		return
	}

	var awardString string
	if claim.Award.Valid {
		awardString = fmt.Sprintf("%.2f", claim.Award.Float64)
	} else {
		awardString = ""
	}
	respondWithJson(w, http.StatusCreated, Response{
		Claims: Claims{
			ID:              claim.ID,
			CustomerID:      claim.CustomerID,
			AssignedAgentID: claim.AgentID,
			ClaimType:       claim.ClaimType,
			CreatedAt:       claim.CreatedAt,
			UpdatedAt:       claim.UpdatedAt,
			CurrentStatus:   string(claim.CurrentStatus),
			AwardAmount:     awardString,
		},
	})
}
