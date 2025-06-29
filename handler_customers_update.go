package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Home       string `json:"home"`
		PolicyType string `json:"policy_type"`
	}

	type Response struct {
		Customer
	}

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
	agentID, err := auth.ValidateJWT(token, cfg.JWT_TOKEN)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	agent, err := cfg.db.GetAgentByID(r.Context(), agentID)
	if err != nil || agentID != agent.ID {
		log.Println(agentID)
		log.Println(agent.ID)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, only agents can update customers", err)
		return
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	cust, err := cfg.db.UpdateCustomer(r.Context(), database.UpdateCustomerParams{
		ID:         custID,
		FirstName:  params.FirstName,
		LastName:   params.LastName,
		Email:      params.Email,
		Phone:      params.Phone,
		Home:       params.Home,
		PolicyType: params.PolicyType,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update customer", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Customer: Customer{
			ID:         custID,
			FirstName:  cust.FirstName,
			LastName:   cust.LastName,
			Email:      cust.Email,
			Phone:      cust.Phone,
			Home:       cust.Home,
			PolicyType: cust.PolicyType,
		},
	})
}
