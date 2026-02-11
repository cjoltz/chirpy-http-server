package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cjoltz/chirpy-http-server/internal/auth"
	"github.com/cjoltz/chirpy-http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Decode Request
	type loginRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process login request", err)
		return
	}

	// Get user details from db
	u, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, u.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong on our end", err)
		return
	}
	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Create access token
	defaultExpiration := 1 * 60 * 60 // 1 hour access token duration
	token, err := auth.MakeJWT(u.ID, cfg.jwtSecret, time.Duration(defaultExpiration)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong: Failed to generate token", err)
		return
	}

	// Create Refresh token and store in database
	refreshTokenDuration := 60 * 24 // 60 day duration in hours
	refreshToken, _ := auth.MakeRefreshToken()
	cfg.db.CreateToken(r.Context(), database.CreateTokenParams{
		Token:     refreshToken,
		UserID:    u.ID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(refreshTokenDuration) * time.Hour),
	})

	type loginResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	respondWithJSON(w, http.StatusOK, loginResponse{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
