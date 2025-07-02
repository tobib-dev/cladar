package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	tokenString, err := auth.MakeJWT(user.ID, cfg.JWT_TOKEN, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT token", err)
		return
	}
	token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	rToken, err := cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440),
	})

	respondWithJson(w, http.StatusOK, Response{
		User: User{
			ID:       user.ID,
			Email:    user.Email,
			UserRole: UserType(user.UserRole),
		},
		Token:        tokenString,
		RefreshToken: rToken.Token,
	})
}
