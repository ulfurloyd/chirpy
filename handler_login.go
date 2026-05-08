package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ulfurloyd/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email                 string      `json:"email"`
		Password              string      `json:"password"`
		ExpiresInSeconds      int         `json:"expires_in_seconds"`
	}
	params := parameters {}

	type response struct {
		User
		Token string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expiresIn := time.Duration(params.ExpiresInSeconds) * time.Second
	if params.ExpiresInSeconds <= 0 || expiresIn > time.Hour {
		expiresIn = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
		Token: token,
	})
}
