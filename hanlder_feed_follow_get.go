package main

import (
	"net/http"

	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedsFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedsFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}
