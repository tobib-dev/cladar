package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetAllManagers(w http.ResponseWriter, r *http.Request) {
	mgrs, err := cfg.db.GetAllManagers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch managers", err)
		return
	}

	respondWithJson(w, http.StatusOK, mgrs)
}
