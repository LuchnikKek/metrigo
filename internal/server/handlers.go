package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/go-chi/chi/v5"
)

type MetricsHandler struct {
	st storage.Storage
}

func NewMetricsHandler(st storage.Storage) *MetricsHandler {
	return &MetricsHandler{st: st}
}

func (h *MetricsHandler) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	m, err := models.ParseMetric(
		chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value"),
	)
	if err != nil {
		log.Println("Bad Request: ", err)
		http.Error(w, "invalid metric data", http.StatusBadRequest)
		return
	}

	if err := h.st.SaveMetric(m); err != nil {
		log.Println("Failed to save metric:", err)
		http.Error(w, "Failed to save metric", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricsHandler) ReadMetricHandler(w http.ResponseWriter, r *http.Request) {
	m, err := h.st.GetMetricByName(chi.URLParam(r, "name"))

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
	w.Write([]byte(fmt.Sprint(m.GetValue())))
}


func (h *MetricsHandler) ReadAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	ms := h.st.GetMetrics()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ms); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
