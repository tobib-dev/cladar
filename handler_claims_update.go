package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerChangeAssignedAgent(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		ID              string `json:"id"`
		AssignedAgentID string `json:"assigned_agent_id"`
	}

	type Response struct {
		Claims
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	claimID, err := uuid.Parse(params.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse claim ID", err)
		return
	}
	agentID, err := uuid.Parse(params.AssignedAgentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse agent ID", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Malformed header; Couldn't retrieve token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized,
			"Token expired or revoked; Please generate new token", err)
		return
	}

	currentClaim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}
	if currentClaim.AgentID == agentID {
		respondWithError(w, http.StatusMethodNotAllowed, "Agent already assigned to this claim", err)
		return
	}

	updatedClaim, err := cfg.db.ChangeAssignedAgent(r.Context(), database.ChangeAssignedAgentParams{
		ID:      claimID,
		AgentID: agentID,
	})

	awardString := ""
	if updatedClaim.Award.Valid {
		awardString = fmt.Sprintf("%.2f", updatedClaim.Award.Float64)
	}

	respondWithJson(w, http.StatusOK, Response{
		Claims: Claims{
			ID:              updatedClaim.ID,
			CustomerID:      updatedClaim.CustomerID,
			AssignedAgentID: updatedClaim.AgentID,
			ClaimType:       updatedClaim.ClaimType,
			CreatedAt:       updatedClaim.UpdatedAt,
			UpdatedAt:       updatedClaim.UpdatedAt,
			CurrentStatus:   string(updatedClaim.CurrentStatus),
			AwardAmount:     awardString,
		},
	})
}
