package main

import (
	"net/http"
	"sort"
	"strings"

	"github.com/Ehs-Colin/chirpy/internal/database"
	"github.com/google/uuid"
)

func authorIdFromRequest(r *http.Request) (uuid.UUID, error) {
	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString == "" {
		return uuid.Nil, nil
	}
	authorId, err := uuid.Parse(authorIdString)
	if err != nil {
		return uuid.Nil, err
	}
	return authorId, nil
}

func getSortOrder(r *http.Request) string {
	sortOrder := strings.ToLower(r.URL.Query().Get("sort"))
	if sortOrder != "asc" && sortOrder != "desc" {
		return "asc"
	}
	return sortOrder
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	// 0> Check if optional author ID was made in request
	authorId, err := authorIdFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author id", err)
		return
	}
	// 1> Get select chirps.  If author is provided, only select that authors chirps
	var dbChirps []database.Chirp
	if authorId == uuid.Nil {
		dbChirps, err = cfg.db.GetChirps(r.Context())
	} else {
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorId)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps", err)
		return
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
	// 4> Sort based on sort order
	sortOrder := getSortOrder(r)
	sort.Slice(returnChirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return returnChirps[i].CreatedAt.After(returnChirps[j].CreatedAt)
		}
		return returnChirps[i].CreatedAt.Before(returnChirps[j].CreatedAt)
	})
	// 5> Return Chirps
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
