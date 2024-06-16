package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	type response struct {
		Feed       database.Feed        `json:"feed"`
		FeedFollow database.FeedsFollow `json:"feed_folow"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	}

	dbFeed, err := apiCfg.DB.CreateFeed(r.Context(), feed)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	feedFollow := database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	dbFeedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), feedFollow)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	resp := response{
		Feed:       dbFeed,
		FeedFollow: dbFeedFollow,
	}

	respondWithJSON(w, http.StatusOK, resp)

}
