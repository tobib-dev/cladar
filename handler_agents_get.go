package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
)

func (cfg *apiConfig) handlerGetAllAgents(w http.ResponseWriter, r *http.Request) {
	agentList, err := cfg.db.GetAllAgents(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting agents from database", err)
		return
	}

	agents := []Agent{}
	for _, agent := range agentList {
		agents = append(agents, Agent{
			ID:        agent.ID,
			FirstName: agent.FirstName,
			LastName:  agent.LastName,
			CreatedAt: agent.CreatedAt,
			UpdatedAt: agent.UpdatedAt,
			Dept:      agent.Dept,
		})
	}

	respondWithJson(w, http.StatusCreated, agents)
}

func (cfg *apiConfig) handlerGetAgent(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Agent
	}

	agentIDString := r.PathValue("agentID")
	agentID, err := uuid.Parse(agentIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse agentID", err)
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
		respondWithError(w, http.StatusUnauthorized, "Token is expired or revoked, please generate new tokens", err)
		return
	}

	agent, err := cfg.db.GetAgentByID(r.Context(), agentID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve agent", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Agent: Agent{
			ID:        agent.ID,
			FirstName: agent.FirstName,
			LastName:  agent.LastName,
			CreatedAt: agent.CreatedAt,
			UpdatedAt: agent.UpdatedAt,
			Email:     agent.Email,
			Dept:      agent.Dept,
		},
	})
}
