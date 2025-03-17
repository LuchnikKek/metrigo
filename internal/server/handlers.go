package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/go-chi/chi/v5"
)

func UpdateMetricHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := models.NewMetric(
			chi.URLParam(r, "type"),
			chi.URLParam(r, "name"),
			chi.URLParam(r, "value"),
		)
		if err != nil {
			log.Println("Internal Server error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := store.Save(m); err != nil {
			log.Println("Failed to save metric:", err)
			http.Error(w, "Failed to save metric", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func ReadMetricHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := store.Get(chi.URLParam(r, "name"))

		if err != nil {
			if err == storage.ErrMetricNotFound {
				http.Error(w, "Metric not found", http.StatusNotFound)
			} else {
				log.Println("Internal Server error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if string(m.GetType()) != chi.URLParam(r, "type") {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(m.GetValue()))
	}
}

func ReadAllMetricsHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ms := store.GetAll()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(ms); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}
}
