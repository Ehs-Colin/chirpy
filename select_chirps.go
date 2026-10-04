package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerChirpsSelect(w http.ResponseWriter, r *http.Request) {
	returnChirps := []Chirp{}
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps", err)
		return
	}
	for _, dbChirp := range dbChirps {
		returnChirps = append(returnChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, returnChirps)
}
