package main

import (
	"net/http"
	"sort"
	"strconv"
)


func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r*http.Request) {
	dbChrips, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChrips {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	chripIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chripIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		Body: dbChirp.Body,
	})
}