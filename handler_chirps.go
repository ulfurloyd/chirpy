package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
		UserID     uuid.UUID    `json:"user_id"`
	}
	type returnVals struct {
		ID             uuid.UUID     `json:"id"`
		CreatedAt      time.Time     `json:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at"`
		Body           string        `json:"body"`
		UserID         uuid.UUID     `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params:= parameters{}
	err := decoder.Decode((&params))
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
		UserID:    params.UserID,
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type responseType struct {
		ID             uuid.UUID     `json:"id"`
		CreatedAt      time.Time     `json:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at"`
		Body           string        `json:"body"`
		UserID         uuid.UUID     `json:"user_id"`
	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}
	responses := []responseType{}
	for _, chirp := range chirps {
		responses = append(responses, responseType{
			ID:          chirp.ID,
			CreatedAt:   chirp.CreatedAt,
			UpdatedAt:   chirp.UpdatedAt,
			Body:        chirp.Body,
			UserID:      chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, responses)
}
