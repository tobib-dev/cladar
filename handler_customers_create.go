package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/database"
)

type Customer struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Home       string    `json:"home"`
	PolicyType string    `json:"policy_type"`
}

func (cfg *apiConfig) handlerCreateCustomer(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Customer
	}

	type Parameters struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Home       string `json:"home"`
		PolicyType string `json:"policy_type"`
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Provide first and last name, email, phone number, address, and policy type", err)
		return
	}

	cust, err := cfg.db.CreateCustomer(r.Context(), database.CreateCustomerParams{
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		Email:      params.Email,
		Phone:      params.Phone,
		Home:       params.Home,
		PolicyType: params.PolicyType,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create customer", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Response{
		Customer: Customer{
			FirstName:  cust.FirstName,
			LastName:   cust.LastName,
			PolicyType: cust.PolicyType,
		},
	})
}
