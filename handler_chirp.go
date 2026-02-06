package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cjoltz/chirpy-http-server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	// Decoding
	type chirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := chirpRequest{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	log.Println("Chirp Request Decoded Successfully")

	// Validation
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	log.Println("Chirp Request Validated Successfully")

	// post chirp to db
	c, err := cfg.db.PostChirp(r.Context(), database.PostChirpParams{Body: cleaned, UserID: params.UserID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create to chirp", err)
		return
	}
	log.Println("Chirp posted to DB")

	// Send response
	respondWithJSON(w, http.StatusCreated, convertDBChirpToChirp(c))
}

func validateChirp(body string) (string, error) {
	charLimit := 140
	profanities := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// Check Length
	if len(body) > charLimit {
		return "", errors.New("Chirp is too long")
	}
	cleaned := filterProfanity(body, profanities)
	return cleaned, nil
}

func filterProfanity(msg string, profanities map[string]struct{}) string {
	words := strings.Split(msg, " ")
	censor_msg := "****"
	for i, word := range words {
		if _, ok := profanities[strings.ToLower(word)]; ok {
			words[i] = censor_msg
		}
	}
	return strings.Join(words, " ")
}
