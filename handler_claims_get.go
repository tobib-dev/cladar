package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerGetAllClaims(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	if user.RevokedAt.Valid || user.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized,
			"Token expired or revoked; Please generate new token", nil)
		return
	}

	claims, err := cfg.db.GetAllClaims(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claims", err)
		return
	}

	allClaims := GetClaimsArray(claims)

	respondWithJson(w, http.StatusOK, allClaims)
}

func (cfg *apiConfig) handlerGetClaim(w http.ResponseWriter, r *http.Request) {
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
			"Token revoked or expired; Please generate new token", nil)
		return
	}

	claim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}

	awardString := GetAwardString(claim.Award)

	respondWithJson(w, http.StatusOK, Response{
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

func (cfg *apiConfig) handlerGetClaimsByCustomer(w http.ResponseWriter, r *http.Request) {
	custIDString := r.PathValue("custID")
	custID, err := uuid.Parse(custIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse custID", err)
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
			"Token expired or revoked; Please genereate new token", nil)
		return
	}

	claims, err := cfg.db.GetAllClaimsByCust(r.Context(), custID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claims", err)
		return
	}

	allCustClaims := GetClaimsArray(claims)

	respondWithJson(w, http.StatusOK, allCustClaims)
}

func (cfg *apiConfig) handlerGetClaimsByAssignedAgent(w http.ResponseWriter, r *http.Request) {
	agentIDString := r.PathValue("agentID")
	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Malformed header; Couldn't parse agentID", err)
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

	claims, err := cfg.db.GetAllClaimsByAgent(r.Context(), agentID)
	if err != nil {
		respondWithError(w, http.StatusNotFound,
			"Couldn't retrieve claims assigned to agent", err)
		return
	}

	agentClaims := GetClaimsArray(claims)

	respondWithJson(w, http.StatusOK, agentClaims)
}

func (cfg *apiConfig) handlerGetPendingClaims(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,
			"Malformed request; Couldn't retrieve bearer token", err)
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

	claims, err := cfg.db.GetPendingClaims(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No pending claims", err)
		return
	}

	claimsSlice := GetClaimsArray(claims)
	respondWithJson(w, http.StatusOK, claimsSlice)
}

func GetClaimsArray(dbClaims []database.Claim) []Claims {
	claims := make([]Claims, len(dbClaims))

	for i, claim := range dbClaims {
		awardString := GetAwardString(claim.Award)

		claims[i] = Claims{
			ID:              claim.ID,
			CustomerID:      claim.CustomerID,
			AssignedAgentID: claim.AgentID,
			ClaimType:       claim.ClaimType,
			CreatedAt:       claim.CreatedAt,
			UpdatedAt:       claim.UpdatedAt,
			CurrentStatus:   string(claim.CurrentStatus),
			AwardAmount:     awardString,
		}
	}

	return claims
}
