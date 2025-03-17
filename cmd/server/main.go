package main

import (
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/server"
	"github.com/LuchnikKek/metrigo/internal/storage"
)

func main() {
	store := storage.NewInMemoryStorage()

	router := server.MetricsRouter(store)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server is running on port 8080")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
