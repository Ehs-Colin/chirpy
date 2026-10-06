package main

import (
	"fmt"
	"net/http"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/Ehs-Colin/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	//1> Validate
	//2> Get Chirp ID from request
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Invalid chirp ID", err)
		return
	}
	//2> Get user from access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate token", err)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate User", err)
		return
	}
	//3> Find chirp in db
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to retreive chirp", err)
		return
	}
	//4> Verify user owns the chirp (This may be redundant with SQL query update?!)
	fmt.Printf("dbChirp Id: %s\n", dbChirp.UserID.String())
	fmt.Printf("userId: %s\n", userId.String())
	if dbChirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Unauthorized deletion attempt", nil)
		return
	}
	//5> Delete the chirp from the database
	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		UserID: userId,
		ID:     dbChirp.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized deletion attempt", nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
