package main

import (
	"encoding/json"
	"net/http"

	"github.com/tobib-dev/cladar/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		User
	}

	params := Parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
		return
	}
	err = auth.VerifyPassword(user.Pswd, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "You don't have access to this account", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		User: User{
			ID:       user.ID,
			Email:    user.Email,
			UserRole: UserType(user.UserRole),
		},
	})
}
