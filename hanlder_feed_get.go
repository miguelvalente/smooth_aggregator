package main

import (
	"net/http"
)

func (apiCfg *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, _ := apiCfg.DB.GetFeeds(r.Context())

	respondWithJSON(w, http.StatusOK, feeds)
}
