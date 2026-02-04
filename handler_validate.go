package main

import (
	"encoding/json"
	"net/http"
	"strings"
)


func handlerChirpValidation(w http.ResponseWriter, r *http.Request) {
	charLimit := 140
	profanities := map[string]struct{} {
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
	}

	type chirpBody struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
	}


	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	if err := decoder.Decode(&chirp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(chirp.Body) > charLimit {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned: filterProfanity(chirp.Body, profanities),
	})
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
