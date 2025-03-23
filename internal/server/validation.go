package server

import (
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/go-chi/chi/v5"
)

func ValidateMetricType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "type")

		// Проверка на пустую строку
		if mType == "" {
			http.Error(w, "metric type is required", http.StatusBadRequest)
			return
		}

		// Проверка на валидность значения
		if !models.IsValidMetricType(mType) {
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateMetricName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mName := chi.URLParam(r, "name")

		// Проверка на пустую строку
		if mName == "" {
			http.Error(w, "metric name is required", http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateMetricValue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "type")
		mName := chi.URLParam(r, "name")
		mValue := chi.URLParam(r, "value")

		// Проверка на пустую строку
		if mValue == "" {
			http.Error(w, "metric value is required", http.StatusBadRequest)
			return
		}

		// Проверка что метрика парсится
		_, err := models.ParseMetric(mType, mName, mValue)
		if err == models.ErrInvalidMetricType {
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "invalid metric", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
