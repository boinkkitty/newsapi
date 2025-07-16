package main

import (
	"github.com/boinkkitty/newsapi/internal/router"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()
	if err := http.ListenAndServe(":5002", r); err != nil {
		log.Fatal("Failed to start server", err)
	}

}
