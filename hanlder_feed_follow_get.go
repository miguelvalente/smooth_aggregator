package main

import (
	"fmt"
	"net/http"

	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedsFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {

	fmt.Println("waaaa")
	feedFollows, err := apiCfg.DB.GetFeedsFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	respondWithJSON(w, http.StatusOK, feedFollows)
}
