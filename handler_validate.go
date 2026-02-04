package main

import (
	"encoding/json"
	"net/http"
)

func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	charLimit := 140

	type chirpBody struct {
		Body string `json:"body"`
	}

	type validResponse struct {
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	// TODO: Consider using utf8.RuneCountInString
	if len(chirp.Body) > charLimit {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, true)
}
