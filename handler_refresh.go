package main

import (
	"net/http"
	"time"

	"github.com/tobib-dev/cladar/internal/auth"
	"github.com/tobib-dev/cladar/internal/database"
)

func (cfg *apiConfig) handlerRefreshTokens(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}

	rToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	if rToken.ExpiresAt.Before(time.Now()) || rToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Tokens have expired or revoked", err)
		return
	}

	freshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't refresh token", err)
		return
	}

	_, err = cfg.db.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:   rToken.Token,
		Token_2: freshToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(rToken.UserID, cfg.JWT_TOKEN, time.Minute*15)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate new access token", err)
		return
	}

	respondWithJson(w, http.StatusOK, Response{
		AccessToken:  accessToken,
		RefreshToken: freshToken,
	})
}

func (cfg *apiConfig) handlerRevokeTokens(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get bearer token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate bearer token", err)
		return
	}

	_, err = cfg.db.RevokeToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
