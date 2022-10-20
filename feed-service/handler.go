package main

import (
	"cqrs/events"
	"cqrs/means"
	"cqrs/models"
	"cqrs/repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/ksuid"
)

type creatFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func handlerCreatedFeed(w http.ResponseWriter, r *http.Request) {
	var req creatFeedRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		means.ErrRes(http.StatusBadRequest, err.Error(), w)
		return
	}

	createdAt := time.Now().UTC()

	id, err := ksuid.NewRandom()
	if err != nil {
		means.ErrRes(http.StatusInternalServerError, err.Error(), w)
		return
	}

	feed := models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   createdAt,
	}

	valid := validator.New()

	if err = valid.Struct(feed); err != nil {
		means.ErrRes(http.StatusBadRequest, err.Error(), w)
		return
	}

	if err := repository.InsertFeed(r.Context(), &feed); err != nil {
		means.ErrRes(http.StatusInternalServerError, err.Error(), w)
	}

	if err := events.PublishCreatedFeed(r.Context(), &feed); err != nil {
		log.Printf("failed to publish created feed event : %v", err)
	}

	means.SuccessRes(http.StatusCreated, "created feed", w)

}
