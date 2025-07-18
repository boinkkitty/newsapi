package handler

import (
	"encoding/json"
	"net/http"

	"github.com/boinkkitty/newsapi/internal/logger"
	"github.com/boinkkitty/newsapi/internal/store"
	"github.com/google/uuid"
)

type NewsStorer interface {
	Create(store.News) (store.News, error)
	FindByID(uuid.UUID) (store.News, error)
	FindAll() ([]store.News, error)
	DeleteByID(uuid.UUID) error
	UpdateByID(body store.News) error
}

func PostNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("post request received")

		var newsRequestBody NewsPostReqBody
		// Parse request from body
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			logger.Error("Failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate fields
		news, err := newsRequestBody.Validate()
		if err != nil {
			logger.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Create news in db
		if _, err := ns.Create(news); err != nil {
			logger.Error("error creating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// To add back response
		w.WriteHeader(http.StatusCreated)
	}
}

func GetAllNews(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("getAll request received")

		// Get all news
		news, err := ns.FindAll()
		if err != nil {
			logger.Error("error getting all news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Create news JSON response
		allNewsResponse := AllNewsResponse{News: news}
		if err := json.NewEncoder(w).Encode(allNewsResponse); err != nil {
			logger.Error("failed to write response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("getByID request received")
		// Reference to router path
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		// Error parsing ID to UUID
		if err != nil {
			logger.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find by UUID
		news, err := ns.FindByID(newsUUID)
		if err != nil {
			logger.Error("news not found", "newsId", newsID)
			w.WriteHeader(http.StatusNotFound)
		}

		// Encode in response
		if err := json.NewEncoder(w).Encode(news); err != nil {
			logger.Error("fail to encode", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func UpdateNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("updateByID request received")

		var newsRequestBody NewsPostReqBody
		// Parse request from body
		if err := json.NewDecoder(r.Body).Decode(&newsRequestBody); err != nil {
			logger.Error("Failed to decode the request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate fields
		news, err := newsRequestBody.Validate()
		if err != nil {
			logger.Error("request validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Update By ID, will be provided in news request body itself
		if err := ns.UpdateByID(news); err != nil {
			logger.Error("error updating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeleteNewsByID(ns NewsStorer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.FromContext(r.Context())
		logger.Info("deleteByID request received")
		// Reference to router path
		newsID := r.PathValue("news_id")
		newsUUID, err := uuid.Parse(newsID)
		// Error parsing ID to UUID
		if err != nil {
			logger.Error("news id not a valid uuid", "newsId", newsID, "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Delete By ID
		if err := ns.DeleteByID(newsUUID); err != nil {
			logger.Error("error updating news", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
