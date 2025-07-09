package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
			"Token expired or revoked; Please generate new token", nil)
		return
	}
	if user.UserRole != database.UserType(UserRoleManager) {
		respondWithError(w, http.StatusForbidden,
			"Only managers have permission to reassign agents", nil)
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

	awardString := GetAwardString(updatedClaim.Award)

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

func (cfg *apiConfig) handlerChangeClaimType(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		ClaimType string `json:"claim_type"`
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

	claimIDString := r.PathValue("claimID")
	claimID, err := uuid.Parse(claimIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse claimID", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't retrieve token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized,
			"Token expired or revoked; Please generate new token", nil)
		return
	}

	currentClaim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}
	if currentClaim.ClaimType == params.ClaimType {
		respondWithError(w, http.StatusMethodNotAllowed,
			"New claim type is the same as current claim type", nil)
		return
	}

	updatedClaim, err := cfg.db.ChangeClaimType(r.Context(), database.ChangeClaimTypeParams{
		ID:        claimID,
		ClaimType: params.ClaimType,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't change claim type", err)
		return
	}

	awardString := GetAwardString(updatedClaim.Award)

	respondWithJson(w, http.StatusOK, Response{
		Claims: Claims{
			ID:              updatedClaim.ID,
			CustomerID:      updatedClaim.CustomerID,
			AssignedAgentID: updatedClaim.AgentID,
			ClaimType:       updatedClaim.ClaimType,
			CreatedAt:       updatedClaim.CreatedAt,
			UpdatedAt:       updatedClaim.UpdatedAt,
			CurrentStatus:   string(updatedClaim.CurrentStatus),
			AwardAmount:     awardString,
		},
	})
}

func (cfg *apiConfig) handlerDeclineClaim(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Claims
	}

	claimIDString := r.PathValue("claimID")
	claimID, err := uuid.Parse(claimIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse claimID", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Malformed header; Couldn't get token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized,
			"Token expired or revoked; Please generate new token", nil)
		return
	}

	currentClaim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}
	if currentClaim.CurrentStatus == database.Status(StatusDeclined) {
		log.Printf("Current status of claim %s: %s", claimID, currentClaim.CurrentStatus)
		respondWithError(w, http.StatusMethodNotAllowed,
			"Cannot change claim status from declined to declined", nil)
		return
	}

	updatedClaim, err := cfg.db.DeclineClaim(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decline claim", err)
		return
	}
	awardString := GetAwardString(updatedClaim.Award)

	respondWithJson(w, http.StatusOK, Response{
		Claims: Claims{
			ID:              updatedClaim.ID,
			CustomerID:      updatedClaim.CustomerID,
			AssignedAgentID: updatedClaim.AgentID,
			ClaimType:       updatedClaim.ClaimType,
			CreatedAt:       updatedClaim.CreatedAt,
			UpdatedAt:       updatedClaim.UpdatedAt,
			CurrentStatus:   string(updatedClaim.CurrentStatus),
			AwardAmount:     awardString,
		},
	})
}

func (cfg *apiConfig) handlerAwardClaim(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		AwardAmount string `json:"award_amount"`
	}

	type Response struct {
		Claims
	}

	claimIDString := r.PathValue("claimID")
	claimID, err := uuid.Parse(claimIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse claimID", err)
		return
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
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
			"Token expired or revoked; Please generate new token", nil)
		return
	}

	currentClaim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}
	if currentClaim.CurrentStatus == database.Status(StatusAwarded) {
		respondWithError(w, http.StatusMethodNotAllowed,
			"Cannot change claim from status awarded to awarded", nil)
		return
	}

	awardSqlFlt, err := GetAwardFloat(params.AwardAmount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't parse award amount to Float64", err)
		return
	}
	updatedClaim, err := cfg.db.ApproveClaim(r.Context(), database.ApproveClaimParams{
		ID:    claimID,
		Award: awardSqlFlt,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't update status to awarded", err)
		return
	}

	awardString := GetAwardString(updatedClaim.Award)
	respondWithJson(w, http.StatusOK, Response{
		Claims: Claims{
			ID:              updatedClaim.ID,
			CustomerID:      updatedClaim.CustomerID,
			AssignedAgentID: updatedClaim.AgentID,
			ClaimType:       updatedClaim.ClaimType,
			CreatedAt:       updatedClaim.CreatedAt,
			UpdatedAt:       updatedClaim.UpdatedAt,
			CurrentStatus:   string(updatedClaim.CurrentStatus),
			AwardAmount:     awardString,
		},
	})
}

func (cfg *apiConfig) handlerChangeAwardAmount(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		AwardAmount string `json:"award_amount"`
	}

	type Response struct {
		Claims
	}

	claimIDString := r.PathValue("claimID")
	claimID, err := uuid.Parse(claimIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed route; Couldn't parse claimID", err)
		return
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Malformed header; Couldn't retrieve bearer token", err)
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

	newAwardAmount, err := GetAwardFloat(params.AwardAmount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't parse award amount to float64", err)
		return
	}

	if newAwardAmount.Float64 == currentClaim.Award.Float64 {
		respondWithError(w, http.StatusMethodNotAllowed,
			"New award amount is same as old award amount; award amount must be different", err)
		return
	}

	updatedClaim, err := cfg.db.ChangeAwardAmount(r.Context(), database.ChangeAwardAmountParams{
		ID:    claimID,
		Award: newAwardAmount,
	})

	newAwardString := GetAwardString(updatedClaim.Award)
	respondWithJson(w, http.StatusOK, Response{
		Claims: Claims{
			ID:              updatedClaim.ID,
			CustomerID:      updatedClaim.CustomerID,
			AssignedAgentID: updatedClaim.AgentID,
			ClaimType:       updatedClaim.ClaimType,
			CreatedAt:       updatedClaim.CreatedAt,
			UpdatedAt:       updatedClaim.UpdatedAt,
			CurrentStatus:   string(updatedClaim.CurrentStatus),
			AwardAmount:     newAwardString,
		},
	})
}

func GetAwardString(awardFloat sql.NullFloat64) string {
	if !awardFloat.Valid {
		return ""
	}
	msg := message.NewPrinter(language.English)
	return msg.Sprintf("$%.2f", awardFloat.Float64)
}

func GetAwardFloat(awardString string) (sql.NullFloat64, error) {
	awardFlt, err := strconv.ParseFloat(awardString, 64)
	if err != nil {
		return sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}, err
	}
	return sql.NullFloat64{
		Float64: awardFlt,
		Valid:   true,
	}, nil
}
