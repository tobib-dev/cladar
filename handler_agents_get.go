package main

import (
	"net/http"
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
