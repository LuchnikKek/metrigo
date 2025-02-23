package server

import (
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/middleware"
	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/gorilla/mux"
)

// NewServer настраивает сервер и маршрутизацию
func NewServer(store storage.Storage) *http.Server {
	router := mux.NewRouter()

	// Настраиваем маршруты
	router.HandleFunc("/update/{type}/{name}/{value}", CreateMetricHandler(store)).Methods("POST")

	// Маршрут для отладки
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Metrics server is running"))
	})

	// Оборачиваем маршрутизатор в middleware
	loggedRouter := middleware.LoggingMiddleware(router)

	// Создаём HTTP-сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedRouter,
	}

	log.Println("Server is running on port 8080")
	return server
}
