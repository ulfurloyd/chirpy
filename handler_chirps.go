package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ulfurloyd/chirpy.git/internal/auth"
	"github.com/ulfurloyd/chirpy.git/internal/database"
)

func getCleanedBody(body string) string {
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}

	words := strings.Split(body, " ")
	for idx, word := range words {
		if _, ok := profanities[strings.ToLower(word)]; ok {
			words[idx] = "****"
		}
	}
	cleanedWords := strings.Join(words, " ")
	return cleanedWords
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string       `json:"body"`
	}
	type returnVals struct {
		ID             uuid.UUID     `json:"id"`
		CreatedAt      time.Time     `json:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at"`
		Body           string        `json:"body"`
		UserID         uuid.UUID     `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid auth header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params:= parameters{}
	err = decoder.Decode((&params))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:      getCleanedBody(params.Body),
		UserID:    userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:          chirp.ID,
		CreatedAt:   chirp.CreatedAt,
		UpdatedAt:   chirp.UpdatedAt,
		Body:        chirp.Body,
		UserID:      chirp.UserID,
	})
}
