package main

import (
	"net/http"

	"github.com/ulfurloyd/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid auth header", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't complete action", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
