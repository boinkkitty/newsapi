package handler

import (
	"context"
	"encoding/json"
	"github.com/boinkkitty/newsapi/internal/news"
	"net/http"

	"github.com/boinkkitty/newsapi/internal/logger"
	"github.com/google/uuid"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler.go -package=mockshandler

type NewsStorer interface {
	Create(context.Context, *news.Record) (*news.Record, error)
	FindByID(context.Context, uuid.UUID) (*news.Record, error)
	FindAll(context.Context) ([]*news.Record, error)
	DeleteByID(context.Context, uuid.UUID) error
	UpdateByID(context.Context, uuid.UUID, *news.Record) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("post request received")

		var newsRequestBody NewsPostReqBody
		// Parse request from body
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("Failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate fields
		n, err := newsRequestBody.Validate()
		if err != nil {
			log.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Create news in db
		if _, err := ns.Create(ctx, n); err != nil {
			log.Error("error creating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// To add back response
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("getAll request received")

		// Get all news
		n, err := ns.FindAll(ctx)
		if err != nil {
			log.Error("error getting all news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Create news JSON response
		allNewsResponse := AllNewsResponse{News: n}
		if err := json.NewEncoder(w).Encode(allNewsResponse); err != nil {
			log.Error("failed to write response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("getByID request received")
		// Reference to router path
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		// Error parsing ID to UUID
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find by UUID
		n, err := ns.FindByID(ctx, newsUUID)
		if err != nil {
			log.Error("news not found", "newsId", newsID)
			w.WriteHeader(http.StatusNotFound)
		}

		// Encode in response
		if err := json.NewEncoder(w).Encode(n); err != nil {
			log.Error("fail to encode", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func UpdateNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("updateByID request received")

		var newsRequestBody NewsPostReqBody
		// Parse request from body
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			log.Error("Failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate fields
		n, err := newsRequestBody.Validate()
		if err != nil {
			log.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Update By ID, will be provided in news request body itself
		if err := ns.UpdateByID(ctx, n.ID, n); err != nil {
			log.Error("error updating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)
		log.Info("deleteByID request received")
		// Reference to router path
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		// Error parsing ID to UUID
		if err != nil {
			log.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Delete By ID
		if err := ns.DeleteByID(ctx, newsUUID); err != nil {
			log.Error("error updating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
