package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerDeleteClaim(w http.ResponseWriter, r *http.Request) {
	claimIDString := r.PathValue("claimID")
	claimID, err := uuid.Parse(claimIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse claimID", err)
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
			"Only user with manager access can delete claims", nil)
		return
	}

	err = cfg.db.DeleteClaim(r.Context(), claimID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete claim", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
