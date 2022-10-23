package main

import (
	"context"
	"cqrs/events"
	"cqrs/means"
	"cqrs/models"
	"cqrs/repository"
	"cqrs/search"
	"encoding/json"
	"log"
	"net/http"
)

func onCreatedFeed(m events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Printf("failed to index feed: %v", err)
	}
}

func handlerListFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	feeds, err := repository.ListFeeds(ctx)

	if err != nil {
		means.ErrRes(http.StatusInternalServerError, err.Error(), w)
	}

	//means.SuccessRes(http.StatusOK, , w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func handlerSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	// - /search?q=value
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		means.ErrRes(http.StatusBadRequest, "query is required", w)
		return
	}

	feeds, err := search.SearchFeed(ctx, query)

	if err != nil {
		means.ErrRes(http.StatusInternalServerError, err.Error(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)

}
