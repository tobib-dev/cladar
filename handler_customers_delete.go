package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
)

func (cfg *apiConfig) handlerDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	custIDString := r.PathValue("custID")
	custID, err := uuid.Parse(custIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't retrieve customer id", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.JWT_TOKEN)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	_, err = cfg.db.GetUserById(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusForbidden, "Access Denied, user does not exist or is unauthorized", err)
		} else {
			respondWithError(w, http.StatusInternalServerError, "Couldn't verify user", err)
		}
		return
	}

	_, err = cfg.db.GetCustomerByID(r.Context(), custID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find customer", err)
		return
	}

	err = cfg.db.DeleteCustomer(r.Context(), custID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to delete customer", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
