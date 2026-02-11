package main

import (
	"net/http"
	"time"

	"github.com/cjoltz/chirpy-http-server/internal/auth"
	"github.com/cjoltz/chirpy-http-server/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	var requestedToken database.RefreshToken
	// No body in request, check header for token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not in request", err)
	}
	// Get user and token details from db
	requestedToken, err = cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token does not exist", err)
		return
	}

	// Create a new ACCESS token
	accessToken, err := auth.MakeJWT(requestedToken.UserID, cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server errror creating access token", err)
		return
	}

	type resp struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, resp{Token: accessToken})
}
