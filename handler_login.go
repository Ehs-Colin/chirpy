package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ehs-Colin/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode create user parameters", err)
		return
	}

	databaseUser, err := cfg.db.SelectUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Incorrect email or password", err)
		return
	}
	validPassword, err := auth.CheckPasswordHash(params.Password, databaseUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Incorrect email or password", err)
		return
	}
	if validPassword {
		userReturn := userReturn{
			ID:        databaseUser.ID,
			CreatedAt: databaseUser.CreatedAt,
			UpdatedAt: databaseUser.UpdatedAt,
			Email:     databaseUser.Email,
		}
		respondWithJSON(w, http.StatusOK, userReturn)
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
	}
}
