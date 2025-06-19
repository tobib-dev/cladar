package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllCustomers(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Customer
	}

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
