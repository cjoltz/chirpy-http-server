package main

import (
	"net/http"

	"github.com/cjoltz/chirpy-http-server/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not in request", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
		respondWithError(w, http.StatusNotFound, "Failed to revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
