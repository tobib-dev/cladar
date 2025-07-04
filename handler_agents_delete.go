package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerDeleteAgent(w http.ResponseWriter, r *http.Request) {
	agentIDString := r.PathValue("agentID")
	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get agentID", err)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed header", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.RevokedAt.Valid || user.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized,
			"Token expired or revoked, please generate new token", nil)
		return
	}

	// Validate user account, user must be manager role to delete agent
	if user.UserRole != database.UserType(UserRoleManager) {
		respondWithError(w, http.StatusForbidden, "Only managers can delete agents", err)
		return
	}

	// Delete agent's user record then delete agent record
	err = cfg.db.DeleteUserByRoleId(r.Context(), agentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete agent's user account", err)
		return
	}
	err = cfg.db.DeleteAgent(r.Context(), agentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete agent", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
