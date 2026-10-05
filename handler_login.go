package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/Ehs-Colin/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode create user parameters", err)
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
	if !validPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	accessToken, err := auth.MakeJWT(databaseUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create access JWT", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	tokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    databaseUser.ID,
		ExpiresAt: time.Now().Add(time.Duration(time.Hour * 24 * 60)),
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), tokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to save refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:           databaseUser.ID,
		CreatedAt:    databaseUser.CreatedAt,
		UpdatedAt:    databaseUser.UpdatedAt,
		Email:        databaseUser.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
