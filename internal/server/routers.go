package server

import (
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricsRouter(store storage.Storage) chi.Router {
	r := chi.NewRouter()
	h := NewMetricsHandler(store)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// GET /
	r.Get("/", h.ReadAllMetricsHandler)

	r.Route("/value", func(r chi.Router) {
		// GET /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
		r.With(ValidateMetricType, ValidateMetricName).
			Get("/{type}/{name}", h.ReadMetricHandler)

		r.Get("/{type}/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "metric name is required", http.StatusBadRequest)
		})
	})

	r.Route("/update", func(r chi.Router) {
		// POST /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		r.With(ValidateMetricType, ValidateMetricName, ValidateMetricValue).
			Post("/{type}/{name}/{value}", h.UpdateMetricHandler)

		r.Post("/{type}/{name}/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "metric value is required", http.StatusBadRequest)
		})
	})

	return r
}
