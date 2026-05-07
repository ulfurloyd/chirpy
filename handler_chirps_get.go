package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	type responseType struct {
		ID             uuid.UUID     `json:"id"`
		CreatedAt      time.Time     `json:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at"`
		Body           string        `json:"body"`
		UserID         uuid.UUID     `json:"user_id"`
	}
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
			respondWithError(w, http.StatusBadRequest, "Could not parse ID", err)
			return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, responseType{
		ID: chirp.ID,
		CreatedAt:   chirp.CreatedAt,
		UpdatedAt:   chirp.UpdatedAt,
		Body:        chirp.Body,
		UserID:      chirp.UserID,
	})
}
