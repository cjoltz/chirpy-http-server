package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cjoltz/chirpy-http-server/internal/auth"
	"github.com/cjoltz/chirpy-http-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type requestUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// Decode Request
	decoder := json.NewDecoder(r.Body)
	params := requestUser{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	// Hash Pasword
	pw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	} 
	// Store user in db
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: pw,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user entry", err)
		return
	}
	// Respond with all user info, except the password
	u := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, http.StatusCreated, u)

}
