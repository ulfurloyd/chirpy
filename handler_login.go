package main

import (
	"encoding/json"
	"net/http"

	"github.com/ulfurloyd/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string      `json:"email"`
		Password    string      `json:"password"`
	}
	params := parameters {}

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

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}
