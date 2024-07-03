package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := GetAPIKey(r.Header)
		if err != nil {
			w.WriteHeader(500)
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if user.ApiKey != apiKey {
			respondWithError(w, http.StatusUnauthorized, "bad api key")
			return
		}

		handler(w, r, user)
	})
}
