package router

import (
	"github.com/boinkkitty/newsapi/internal/handler"
	"net/http"
)

func NewRouter(ns handler.NewsStorer) *http.ServeMux {
	// Create router
	router := http.NewServeMux()

	// Create routes
	router.HandleFunc("POST /news", handler.PostNews(ns))
	router.HandleFunc("GET /news", handler.GetAllNews(ns))
	router.HandleFunc("GET /news/{news_id}", handler.GetNewsByID(ns))
	router.HandleFunc("PUT /news/{news_id}", handler.UpdateNewsByID(ns))
	router.HandleFunc("DELETE /news/{news_id}", handler.DeleteNewsByID(ns))

	return router
}
