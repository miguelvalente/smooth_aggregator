package main

import (
	"net/http"

	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}
