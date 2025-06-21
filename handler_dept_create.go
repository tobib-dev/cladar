package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Dept struct {
	ID       uuid.UUID `json:"id"`
	DeptName string    `json:"dept_name"`
}

func (cfg *apiConfig) handlerCreateDept(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		DeptName string `json:"dept_name"`
	}

	type Response struct {
		Dept
	}
	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	dept, err := cfg.db.CreateDept(r.Context(), params.DeptName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create department", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Response{
		Dept: Dept{
			ID:       dept.ID,
			DeptName: dept.DeptName,
		},
	})
}
