package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/database"
)

type Manager struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	DeptID    uuid.UUID `json:"dept_id"`
}

func (cfg *apiConfig) handlerCreateManager(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		DeptName  string `json:"dept"`
		Password  string `json:"password"`
	}

	type Response struct {
		Manager
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	dept, err := cfg.db.GetDeptByID(r.Context(), params.DeptName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find department", err)
		return
	}

	mang, err := cfg.db.CreateManager(r.Context(), database.CreateManagerParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		DeptID:    dept.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create manager", err)
		return
	}

	_, err = cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Pswd:     params.Password,
		UserRole: database.UserTypeManager,
		RoleID:   mang.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user account for manager", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Response{
		Manager: Manager{
			FirstName: mang.FirstName,
			LastName:  mang.LastName,
			Email:     mang.Email,
			DeptID:    mang.DeptID,
		},
	})
}
