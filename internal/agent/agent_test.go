package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricsAgent_Send(t *testing.T) {
	type MetricTestStruct struct {
		Name  string
		Value any
		Type  string
	}
	testCases := []struct {
		name  string
		value MetricTestStruct
		want  string
	}{
		{
			name: "TestSendCounter",
			value: MetricTestStruct{
				Name:  "PollCount",
				Value: 10,
				Type:  "counter",
			},
			want: "/update/counter/PollCount/10",
		},
		{
			name: "TestSendGauge",
			value: MetricTestStruct{
				Name:  "RandomValue",
				Value: 0.1234,
				Type:  "gauge",
			},
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

			client := &http.Client{}
			ag := NewMetricsAgent(client, ts.URL)

			err := ag.Send(tC.value.Name, tC.value.Value, tC.value.Type)
			require.NoError(t, err)
		})
	}
}
