package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

type Agent struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Dept      string    `json:"dept"`
}

func (cfg *apiConfig) handlerCreateAgent(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Dept      string `json:"dept"`
		Password  string `json:"password"`
	}

	type Response struct {
		Agent
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}
	user, err := cfg.db.GetUserFromToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	if user.UserRole != database.UserType(UserRoleManager) {
		respondWithError(w, http.StatusUnauthorized, "Only managers have permission to create agents", err)
		return
	}
	if user.ExpiresAt.Before(time.Now()) || user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token is expired or revoked, please generate new tokens", err)
		return
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	agent, err := cfg.db.CreateAgent(r.Context(), database.CreateAgentParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Dept:      params.Dept,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create agent", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	_, err = cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    agent.Email,
		Pswd:     hashedPassword,
		UserRole: database.UserTypeAgent,
		RoleID:   agent.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user account", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Agent: Agent{
			FirstName: agent.FirstName,
			LastName:  agent.LastName,
			Email:     agent.Email,
			Dept:      agent.Dept,
		},
	})
}
