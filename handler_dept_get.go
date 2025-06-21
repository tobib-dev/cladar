package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllDepts(w http.ResponseWriter, r *http.Request) {
	depts, err := cfg.db.GetAllDept(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respondWithJson(w, http.StatusOK, depts)
}
