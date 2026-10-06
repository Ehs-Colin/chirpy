package main

import (
	"encoding/json"
	"net/http"

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
		Password string `json:"password"`
	}
	// 1> Decode webhook JSON into usable parameters
	decoder := json.NewDecoder(r.Body)
	params := polkaParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode create user parameters", err)
		return
	}
	// 2> If event is not user.upgrade, return 204
	if params.Event != UPGRADE_EVENT {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// 3> If user does not exist, 404
	_, err = cfg.db.SelectUserById(r.Context(), params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}
	// 4> Update user to chirpy red
	_, err = cfg.db.UpdateUserSetChirpyRed(r.Context(), database.UpdateUserSetChirpyRedParams{
		ID:          params.Data.UserId,
		IsChirpyRed: true,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to set chirpy red", err)
		return
	}
	// 5> If it succeeded, 204.
	w.WriteHeader(http.StatusNoContent)

}
