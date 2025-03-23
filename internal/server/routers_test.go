package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/LuchnikKek/metrigo/internal/storage"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestMetricsRouterUpdateGauge(t *testing.T) {
	st := storage.NewInMemoryStorage()
	ts := httptest.NewServer(MetricsRouter(st))
	defer ts.Close()

	var testTable = []struct {
		method string
		url    string
		want   string
		status int
	}{
		{"POST", "/update/gauge/HeapIdle/5210112", "", http.StatusOK},
		{"POST", "/update/gauge/HeapIdle/5210112.12", "", http.StatusOK},

		{"GET", "/update/gauge/HeapIdle/5210112", "", http.StatusMethodNotAllowed},
		{"PUT", "/update/gauge/HeapIdle/5210112", "", http.StatusMethodNotAllowed},
		{"PATCH", "/update/gauge/HeapIdle/5210112", "", http.StatusMethodNotAllowed},
		{"DELETE", "/update/gauge/HeapIdle/5210112", "", http.StatusMethodNotAllowed},

		{"POST", "/update/gauge//5210112", "metric name is required\n", http.StatusNotFound},
		{"POST", "/update//HeapIdle/5210112", "metric type is required\n", http.StatusBadRequest},
		{"POST", "/update/gauge/HeapIdle/", "metric value is required\n", http.StatusBadRequest},
		{"POST", "/update/gauge/HeapIdle/lalala", "invalid metric\n", http.StatusBadRequest},
		{"POST", "/update/invalid/HeapIdle/5210112", "invalid metric type\n", http.StatusBadRequest},
	}
	for _, v := range testTable { // single test with all requests
		resp, get := testRequest(t, ts, v.method, v.url)
		defer resp.Body.Close()
		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}
}

func TestMetricsRouterUpdateCounter(t *testing.T) {
	st := storage.NewInMemoryStorage()
	ts := httptest.NewServer(MetricsRouter(st))
	defer ts.Close()

	var testTable = []struct {
		method string
		url    string
		want   string
		status int
	}{
		{"POST", "/update/counter/PollCount/1", "", http.StatusOK},
		{"POST", "/update/counter/PollCount/1.12", "invalid metric\n", http.StatusBadRequest},

		{"GET", "/update/counter/PollCount/1", "", http.StatusMethodNotAllowed},
		{"PUT", "/update/counter/PollCount/1", "", http.StatusMethodNotAllowed},
		{"PATCH", "/update/counter/PollCount/1", "", http.StatusMethodNotAllowed},
		{"DELETE", "/update/counter/PollCount/1", "", http.StatusMethodNotAllowed},

		{"POST", "/update/counter//1", "metric name is required\n", http.StatusNotFound},
		{"POST", "/update//PollCount/1", "metric type is required\n", http.StatusBadRequest},
		{"POST", "/update/counter/PollCount/", "metric value is required\n", http.StatusBadRequest},
		{"POST", "/update/counter/PollCount/lalala", "invalid metric\n", http.StatusBadRequest},
		{"POST", "/update/invalid/PollCount/1", "invalid metric type\n", http.StatusBadRequest},
	}
	for _, v := range testTable { // single test with all requests
		resp, get := testRequest(t, ts, v.method, v.url)
		defer resp.Body.Close()

		assert.Equal(t, v.status, resp.StatusCode)
		assert.Equal(t, v.want, get)
	}
}

func TestMetricsUpdateOverwriteCounter(t *testing.T) {
	var testTable = []struct {
		first    string
		second   string
		expected int
	}{
		{"1", "2", 3},
		{"10", "4", 14},
		{"0", "100", 100},
		{"44", "0", 44},
	}

	for _, v := range testTable {
		t.Run(fmt.Sprintf("Counter %s updated by %s", v.first, v.second), func(t *testing.T) {
			st := storage.NewInMemoryStorage()
			ts := httptest.NewServer(MetricsRouter(st))
			defer ts.Close()

			resp, _ := testRequest(t, ts, "POST", "/update/counter/PollCount/"+v.first)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()

			resp, _ = testRequest(t, ts, "POST", "/update/counter/PollCount/"+v.second)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()

			value, err := st.GetMetricByName("PollCount")
			require.NoError(t, err)

			assert.Equal(t, value.GetValue(), v.expected)
		})
	}
}

func TestMetricsUpdateOverwriteGauge(t *testing.T) {
	var testTable = []struct {
		first    string
		second   string
		expected float64
	}{
		{"5210", "521", 521},
		{"10", "4", 4},
		{"10", "4.4", 4.4},
		{"10.1", "4", 4},
		{"10.1", "4.4", 4.4},
		{"111", "0", 0},
		{"0", "42.6", 42.6},
	}
	for _, v := range testTable {
		t.Run(fmt.Sprintf("Gauge %s updated by %s", v.first, v.second), func(t *testing.T) {
			st := storage.NewInMemoryStorage()
			ts := httptest.NewServer(MetricsRouter(st))
			defer ts.Close()

			resp, _ := testRequest(t, ts, "POST", "/update/gauge/HeapIdle/"+v.first)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()

			resp, _ = testRequest(t, ts, "POST", "/update/gauge/HeapIdle/"+v.second)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()

			value, err := st.GetMetricByName("HeapIdle")
			require.NoError(t, err)

			assert.Equal(t, value.GetValue(), v.expected)
		})
	}
}

func TestMetricsRouterRead(t *testing.T) {
	st := storage.NewInMemoryStorage()
	st.SaveMetric(models.NewGaugeMetric("HeapInuse", 729088))
	st.SaveMetric(models.NewGaugeMetric("RandomValue", 0.31415926))
	st.SaveMetric(models.NewCounterMetric("PollCount", 2))

	ts := httptest.NewServer(MetricsRouter(st))
	defer ts.Close()

	var testTable = []struct {
		name   string
		method string
		url    string
		want   string
		status int
	}{
		{"gauge int", "GET", "/value/gauge/HeapInuse", "729088", http.StatusOK},
		{"gauge float", "GET", "/value/gauge/RandomValue", "0.31415926", http.StatusOK},
		{"counter int", "GET", "/value/counter/PollCount", "2", http.StatusOK},

		{"metric name not in available", "GET", "/value/gauge/lalala", "Metric not found\n", http.StatusNotFound},
		{"metric name in available, not fetched", "GET", "/value/gauge/StackInuse", "Metric not found\n", http.StatusNotFound},
		{"incorrect metric type", "GET", "/value/counter/RandomValue", "Metric not found\n", http.StatusNotFound},
		{"metric type does not exists", "GET", "/value/lol/RandomValue", "invalid metric type\n", http.StatusBadRequest},

		{"empty type", "GET", "/value//HeapInuse", "metric type is required\n", http.StatusBadRequest},
		{"empty name", "GET", "/value/gauge/", "metric name is required\n", http.StatusBadRequest},

		{"method POST not allowed", "POST", "/value/gauge/HeapInuse", "", http.StatusMethodNotAllowed},
		{"method PUT not allowed", "PUT", "/value/gauge/HeapInuse", "", http.StatusMethodNotAllowed},
		{"method PATCH not allowed", "PATCH", "/value/gauge/HeapInuse", "", http.StatusMethodNotAllowed},
		{"method DELETE not allowed", "DELETE", "/value/gauge/HeapInuse", "", http.StatusMethodNotAllowed},
	}
	for _, v := range testTable {
		t.Run(v.name, func(t *testing.T) {
			resp, get := testRequest(t, ts, v.method, v.url)
			defer resp.Body.Close()
			assert.Equal(t, v.status, resp.StatusCode)
			assert.Equal(t, v.want, get)
		})
	}
}

func TestMetricsRouterReadAll(t *testing.T) {
	st := storage.NewInMemoryStorage()
	expectedData := []models.Metric{
		models.NewCounterMetric("PollCount", 2),
		models.NewGaugeMetric("HeapInuse", 729088),
		models.NewGaugeMetric("RandomValue", 0.31415926),
	}
	for _, mData := range expectedData {
		st.SaveMetric(mData)
	}

	ts := httptest.NewServer(MetricsRouter(st))
	defer ts.Close()

	type MockMetric struct {
		Name  string            `json:"name"`
		Value any           	`json:"value"`
		Type  models.MetricType `json:"type"`
	}
	respData := []MockMetric{}

	client := resty.New()
	resp, err := client.R().
		SetResult(&respData).
		Get(ts.URL + "/")
	require.NoError(t, err)

	assert.Equal(t, resp.StatusCode(), http.StatusOK)
	assert.Equal(t, len(respData), len(expectedData))

	for _, respVal := range respData {
		expVal, err := st.GetMetricByName(respVal.Name)

		assert.NoError(t, err)

		assert.Equal(t, expVal.GetName(), respVal.Name)
		assert.Equal(t, expVal.GetType(), respVal.Type)
		
		if respVal.Type == models.Gauge {
			assert.Equal(t, expVal.GetValue(), respVal.Value)
		}
		if respVal.Type == models.Counter {
			intVal, _ := strconv.Atoi(fmt.Sprint(respVal.Value))
			assert.Equal(t, expVal.GetValue(), intVal)
		}
	}
}
