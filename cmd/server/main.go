package main

import (
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/server"
	"github.com/LuchnikKek/metrigo/internal/storage"
)

func main() {
	InitOptions()
	store := storage.NewInMemoryStorage()

	router := server.MetricsRouter(store)

	srv := &http.Server{
		Addr:    Options.Addr,
		Handler: router,
	}

	log.Println("Server is running on", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
