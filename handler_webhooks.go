package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ehs-Colin/chirpy/internal/auth"
	"github.com/Ehs-Colin/chirpy/internal/database"
	"github.com/google/uuid"
)

const UPGRADE_EVENT = "user.upgraded"

func (cfg *apiConfig) handlerPolkaWebHook(w http.ResponseWriter, r *http.Request) {
	type polkaParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	// 0> Verify API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	// 1> Decode webhook JSON into usable parameters
	decoder := json.NewDecoder(r.Body)
	params := polkaParameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode create user parameters", err)
		return
	}
	// 2> If event is not user.upgrade, return 204
	if params.Event != UPGRADE_EVENT {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// 3> Update user to chirpy red.  404 if user not found
	_, err = cfg.db.UpdateUserSetChirpyRed(r.Context(), database.UpdateUserSetChirpyRedParams{
		ID:          params.Data.UserId,
		IsChirpyRed: true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Unable to set chirpy red", err)
		return
	}
	// 4> If it succeeded, 204.
	w.WriteHeader(http.StatusNoContent)
}
