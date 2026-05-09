package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/ulfurloyd/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handlerChirpyRed(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	req := request{}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API key", err)
		return
	}

	if apiKey != cfg.polkaAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse User ID", err)
		return
	}

	_, err = cfg.db.UpdateUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
