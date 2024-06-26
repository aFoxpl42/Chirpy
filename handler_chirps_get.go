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

	authorID := -1
	authorIDstring := r.URL.Query().Get("author_id")
	if authorIDstring != "" {
		authorID, err = strconv.Atoi(authorIDstring)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	sortType := r.URL.Query().Get("sort")
	if sortType == "" {
		sortType = "asc"
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChrips {
		if authorID != -1 && dbChirp.AuthorID != authorID {
			continue
		}

		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}
	if sortType == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

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