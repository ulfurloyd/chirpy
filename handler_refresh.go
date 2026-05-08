package main

import (
	"net/http"
	"time"

	"github.com/ulfurloyd/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid auth header", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing", err)
		return
	}

	JWTToken, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: JWTToken,
	})
}
