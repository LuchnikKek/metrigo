package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LuchnikKek/metrigo/internal/models"
	"github.com/stretchr/testify/require"
)

func TestMetricsAgent_SendMetric(t *testing.T) {
	testCases := []struct {
		name  string
		value models.Metric
		want  string
	}{
		{
			name: "TestSendCounter",
			value: models.NewCounterMetric("PollCount", 10),
			want: "/update/counter/PollCount/10",
		},
		{
			name: "TestSendGauge",
			value: models.NewGaugeMetric("RandomValue", 0.1234),
			want: "/update/gauge/RandomValue/0.1234",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("%s method = \"%s\", want \"%s\"", tC.name, r.Method, http.MethodPost)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				if r.URL.Path != tC.want {
					t.Errorf("%s path = \"%s\", want \"%s\"", tC.name, r.URL.Path, tC.want)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}))
			defer ts.Close()

			ag := NewMetricsAgent(ts.URL, 2*time.Second, 10*time.Second, 5*time.Second)

			err := ag.SendMetric(tC.value)
			require.NoError(t, err)
		})
	}
}
