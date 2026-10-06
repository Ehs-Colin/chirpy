package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/Ehs-Colin/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type userParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}
	//1> Validate access token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate token", err)
		return
	}
	//2> Get user id from access token
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate User", err)
		return
	}
	//3> Get new email and password from request body
	decoder := json.NewDecoder(r.Body)
	params := userParameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode create user parameters", err)
		return
	}
	//4> Hash password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}
	//5> Update email and hashed password in database for given user id
	databaseUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user", err)
		return
	}
	//6> Return 200 and updated user
	respondWithJSON(w, http.StatusOK, response{
		ID:          databaseUser.ID,
		CreatedAt:   databaseUser.CreatedAt,
		UpdatedAt:   databaseUser.UpdatedAt,
		Email:       databaseUser.Email,
		IsChirpyRed: databaseUser.IsChirpyRed,
	})

}
