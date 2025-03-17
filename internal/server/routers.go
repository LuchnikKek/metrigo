package server

import (
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricsRouter(store storage.Storage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// GET /
	r.Get("/", ReadAllMetricsHandler(store))

	r.Route("/value", func(r chi.Router) {
		// GET /value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
		r.With(ValidateMetricType, ValidateMetricName).
			Get("/{type}/{name}", ReadMetricHandler(store))

		r.Get("/{type}/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "metric name is required", http.StatusBadRequest)
		})
	})

	r.Route("/update", func(r chi.Router) {
		// POST /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		r.With(ValidateMetricType, ValidateMetricName, ValidateMetricValue).
			Post("/{type}/{name}/{value}", UpdateMetricHandler(store))

		r.Post("/{type}/{name}/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "metric value is required", http.StatusBadRequest)
		})
	})

	return r
}
