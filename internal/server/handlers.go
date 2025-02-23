package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/gorilla/mux"
)

func CreateMetricHandler(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println("Parsed vars:", vars)

		metricType, okType := vars["type"]
		metricName, okName := vars["name"]
		metricValue, okValue := vars["value"]

		if !okName || metricName == "" {
			log.Println("Error: Metric name is missing")
			http.Error(w, "Metric name is required", http.StatusNotFound)
			return
		}
		if !okType {
			log.Println("Error: Metric type is missing")
			http.Error(w, "Metric type is required", http.StatusBadRequest)
			return
		}
		if !okValue {
			log.Println("Error: Metric value is missing")
			http.Error(w, "Metric value is required", http.StatusBadRequest)
			return
		}

		log.Printf("Processing metric: type=%s, name=%s, value=%s\n", metricType, metricName, metricValue)

		// Определяем тип метрики и парсим значение
		var metric models.Metric
		var err error

		switch metricType {
		case string(models.Gauge):
			var val float64
			val, err = strconv.ParseFloat(metricValue, 64)
			if err == nil {
				metric = &models.GaugeMetric{Name: metricName, Value: val}
			}
		case string(models.Counter):
			var val int64
			val, err = strconv.ParseInt(metricValue, 10, 64)
			if err == nil {
				metric = &models.CounterMetric{Name: metricName, Value: val}
			}
		default:
			log.Println("Error: Invalid metric type:", metricType)
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Println("Error: Invalid metric value:", metricValue)
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		if err := store.Save(metric); err != nil {
			log.Println("Error: Failed to save metric:", err)
			http.Error(w, "Failed to save metric", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
