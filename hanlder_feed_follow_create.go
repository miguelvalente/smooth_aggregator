package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedsFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	_, err = apiCfg.DB.GetFeedsbyID(r.Context(), params.FeedId)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	feedFollows := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	}

	databaseFeedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), feedFollows)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollow)
}
