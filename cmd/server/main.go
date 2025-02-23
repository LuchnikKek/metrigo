package main

import (
	"log"

	"github.com/LuchnikKek/metrigo/internal/server"
	"github.com/LuchnikKek/metrigo/internal/storage"
)

func main() {
	store := storage.NewInMemoryStorage()

	srv := server.NewServer(store)

	log.Println("Server running on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
