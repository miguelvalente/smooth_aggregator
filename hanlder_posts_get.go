package main

import (
	"net/http"
	"strconv"

	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerPostsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.PathValue("limit")
	var limit int
	if limitStr == "" {
		limit = 3
	} else {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
	}

	posts, err := apiCfg.DB.GetPostsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	if len(posts) > limit {
		posts = posts[:limit]
	}

	respondWithJSON(w, http.StatusOK, posts)
}
