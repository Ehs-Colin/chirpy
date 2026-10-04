package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/Ehs-Colin/chirpy/internal/database"
	"github.com/google/uuid"
)

type userParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userReturn struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := userParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode create user parameters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	databaseUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}
	userReturn := userReturn{
		ID:        databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
		Email:     databaseUser.Email,
	}
	respondWithJSON(w, http.StatusCreated, userReturn)
}
