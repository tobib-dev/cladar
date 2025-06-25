package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
)

func (cfg *apiConfig) handlerGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	custList, err := cfg.db.GetAllCustomers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error pulling customers from database", err)
		return
	}

	customers := []Customer{}
	for _, cust := range custList {
		if cust.ID != uuid.Nil {
			customers = append(customers, Customer{
				ID:         cust.ID,
				FirstName:  cust.FirstName,
				LastName:   cust.LastName,
				Email:      cust.Email,
				Home:       cust.Home,
				Phone:      cust.Phone,
				PolicyType: cust.PolicyType,
			})
		}
	}

	respondWithJson(w, http.StatusOK, customers)
}

func (cfg *apiConfig) handlerGetCustomer(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("custID")
	custID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid customer id", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Bearer token", err)
		return
	}
	_, err = auth.ValidateJWT(token, cfg.JWT_TOKEN)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	cust, err := cfg.db.GetCustomerByID(r.Context(), custID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find customer", err)
		return
	}

	respondWithJson(w, http.StatusOK, Customer{
		ID:         cust.ID,
		FirstName:  cust.FirstName,
		LastName:   cust.LastName,
		Email:      cust.Email,
		Phone:      cust.Phone,
		Home:       cust.Home,
		PolicyType: cust.PolicyType,
	})
}
