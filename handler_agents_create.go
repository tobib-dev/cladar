package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/database"
)

type Agent struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Dept      string    `json:"dept"`
}

func (cfg *apiConfig) handlerCreateAgent(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Dept      string `json:"dept"`
	}

	type Response struct {
		Agent
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	agent, err := cfg.db.CreateAgent(r.Context(), database.CreateAgentParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Dept:      params.Dept,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create agent", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Agent: Agent{
			FirstName: agent.FirstName,
			LastName:  agent.LastName,
			Dept:      agent.Dept,
		},
	})
}
