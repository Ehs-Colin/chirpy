package main

import (
	"net/http"

	"github.com/Ehs-Colin/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error
	authorId := r.URL.Query().Get("author_id")

	// 1> Get select chirps.  If author is provided, only select that authors chirps
	if authorId == "" {
		dbChirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps", err)
			return
		}
	} else {
		authorUUID, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps", err)
			return
		}
	}
	// 3> Add selected Chirps to Chirp struct list
	returnChirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		returnChirps = append(returnChirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}
	// 4> Return Chirps
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
	respondWithJSON(w, http.StatusOK, returnChirp)
}
