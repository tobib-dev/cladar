package main

import (
	"net/http"
	"time"

	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}

	rToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	if rToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Tokens have expired", err)
		return
	}

	freshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't refresh token", err)
		return
	}

	updatedToken, err := cfg.db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:   rToken.Token,
		Token_2: freshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update refresh token", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		Token: updatedToken.Token,
	})
}
