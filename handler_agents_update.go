package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerUpdateAgents(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Agent
	}

	type Parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Dept      string `json:"dept"`
	}

	agentIDString := r.PathValue("agentID")
	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse agentID", err)
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
		respondWithError(w, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}

	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token is expired or has been revoked, please generate new tokens", err)
		return
	}

	isManager := user.UserRole == database.UserTypeManager

	if !isManager && user.RoleID != agentID {
		/*
			 * caller must be a manager or agent that owns the user profile
				* For instance, say agent with agentID 002 wants to change his/her
				* profile they will be able to but agent with agentID 003 cannot
				* change agent 002's profile. Only managers have permission to
				* update others account
		*/
		respondWithError(w, http.StatusForbidden,
			"Only managers and owners can update agent account", nil)
		return
	}

	agent, err := cfg.db.GetAgentByID(r.Context(), agentID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Agent not found", err)
		return
	}

	fName := agent.FirstName
	if params.FirstName != "" {
		fName = params.FirstName
	}
	lName := agent.LastName
	if params.LastName != "" {
		lName = params.LastName
	}
	email := agent.Email
	if params.Email != "" {
		email = params.Email
	}
	dept := agent.Dept
	if params.Dept != "" {
		dept = params.Dept
	}

	dbAgent, err := cfg.db.UpdateAgent(r.Context(), database.UpdateAgentParams{
		ID:        agentID,
		FirstName: fName,
		LastName:  lName,
		Email:     email,
		Dept:      dept,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update agent", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Agent: Agent{
			ID:        agent.ID,
			FirstName: dbAgent.FirstName,
			LastName:  dbAgent.LastName,
			CreatedAt: dbAgent.CreatedAt,
			UpdatedAt: dbAgent.UpdatedAt,
			Email:     dbAgent.Email,
			Dept:      dbAgent.Dept,
		},
	})
}
