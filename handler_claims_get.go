package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
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
			"Token expired or revoked; Please generate new token", err)
		return
	}

	claims, err := cfg.db.GetAllClaims(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claims", err)
		return
	}

	allClaims := make([]Claims, len(claims))
	for i, claim := range claims {
		var awardString string
		if claim.Award.Valid {
			awardString = fmt.Sprintf("%.2f", claim.Award.Float64)
		} else {
			awardString = ""
		}

		allClaims[i] = Claims{
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
			"Token revoked or expired; Please generate new token", err)
		return
	}

	claim, err := cfg.db.GetClaimByID(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claim", err)
		return
	}

	awardString := ""
	if claim.Award.Valid {
		awardString = fmt.Sprintf("%.2f", claim.Award.Float64)
	}

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
			"Token expired or revoked; Please genereate new token", err)
		return
	}

	claims, err := cfg.db.GetAllClaimsByCust(r.Context(), custID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve claims", err)
		return
	}

	allCustClaims := make([]Claims, len(claims))
	for i, claim := range claims {
		awardString := ""
		if claim.Award.Valid {
			awardString = fmt.Sprintf("%.2f", claim.Award.Float64)
		}
		allCustClaims[i] = Claims{
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

	respondWithJson(w, http.StatusOK, allCustClaims)
}
