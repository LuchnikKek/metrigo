package server

import (
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/gorilla/mux"
)

// http://<server>/update/{type}/{name}/{value}
func CreateMetricHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Parsed vars:", vars)

		req, err := ValidateMetricVars(vars)
		if err != nil {
			switch err {
			case ErrMissingMetricName:
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			case ErrMissingMetricType, ErrMissingMetricValue:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			default:
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
		}

		metric, err := ParseMetric(req)
		if err != nil {
			switch err {
			case ErrInvalidMetricValue:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case ErrInvalidMetricType:
				http.Error(w, err.Error(), http.StatusBadRequest)
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		if err := store.Save(metric); err != nil {
			log.Println("Failed to save metric:", err)
			http.Error(w, "Failed to save metric", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
