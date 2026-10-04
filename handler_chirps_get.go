package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsSelectAll(w http.ResponseWriter, r *http.Request) {
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

func (cfg *apiConfig) handlerChirpsSelectById(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to retreive chirp", err)
		return
	}
	returnChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusMethodNotAllowed, returnChirp)
}
