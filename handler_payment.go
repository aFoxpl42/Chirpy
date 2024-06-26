package main

import (
	"encoding/json"
	"net/http"

	"github.com/aFoxpl42/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPayment(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find ApiKey")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong ApiKey")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := cfg.DB.GetUserByID(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user with provided ID")
		return
	}

	err = cfg.DB.UpdateUserRed(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	}