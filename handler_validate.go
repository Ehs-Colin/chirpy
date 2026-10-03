package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleaned := sanitizeChirp((params.Body))
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

func sanitizeChirp(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	var result []string
	for _, word := range strings.Split(body, " ") {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			result = append(result, "****")
		} else {
			result = append(result, word)
		}
	}
	return strings.Join(result, " ")
}
