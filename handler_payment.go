package main

import (
	"net/http"
	"encoding/json"

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