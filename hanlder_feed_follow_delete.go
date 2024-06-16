package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedsFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID := uuid.MustParse(r.PathValue("feedFollowID"))

	// _, err := apiCfg.DB.GetFeedsbyID(r.Context(), feedFollowID)
	// if err != nil {
	// 	respondWithError(w, http.StatusForbidden, err.Error())
	// }

	err := apiCfg.DB.DeleteUserFeedFollows(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	respondWithJSON(w, http.StatusOK, nil)
}
