package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/google/uuid"
)

type loginParameters struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type loginReturn struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := loginParameters{}
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
	if !validPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	expiresInSeconds := params.ExpiresInSeconds
	if expiresInSeconds <= 0 || expiresInSeconds > 3600 {
		expiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(databaseUser.ID, cfg.secret, time.Duration(expiresInSeconds*int(time.Second)))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	userReturn := loginReturn{
		ID:        databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email:     databaseUser.Email,
		Token:     token,
	}
	respondWithJSON(w, http.StatusOK, userReturn)
}
