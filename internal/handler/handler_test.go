package handler

import (
	"github.com/Kuart/metric-collector/internal/storage/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "case 1: unknown metric type",
			request: "/update/unknown/testGauge/1.1",
			want: want{
				statusCode: 501,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			storage := storage.New()
			metricHandler := NewHandler(storage)
			NewRouter(metricHandler)

			h := http.HandlerFunc(metricHandler.Update)
			h.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func TestCounterHandler(t *testing.T) {
	pattern := "/{type}/{name}/{value}"
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "case 1: unknown metric value",
			request: "/update/counter/testGauge/",
			want: want{
				statusCode: 404,
			},
		},
		{
			name:    "case 2: metric value wrong format",
			request: "/update/counter/testGauge/1.1",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:    "case 3: metric value ok",
			request: "/update/counter/testGauge/1",
			want: want{
				statusCode: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			storage := storage.New()
			metricHandler := NewHandler(storage)
			r := NewRouter(metricHandler)

			r.HandleFunc(pattern, metricHandler.Counter)
			r.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}

func TestGaugeHandler(t *testing.T) {
	pattern := "/{type}/{name}/{value}"
	type want struct {
		statusCode int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "case 1: unknown metric value",
			request: "/update/gauge/testGauge/",
			want: want{
				statusCode: 404,
			},
		},
		{
			name:    "case 2: metric value wrong format",
			request: "/update/gauge/testGauge/1s",
			want: want{
				statusCode: 400,
			},
		},
		{
			name:    "case 3: metric value ok",
			request: "/update/gauge/testGauge/1.1",
			want: want{
				statusCode: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			storage := storage.New()
			metricHandler := NewHandler(storage)
			r := NewRouter(metricHandler)

			r.HandleFunc(pattern, metricHandler.Gauge)
			r.ServeHTTP(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			defer result.Body.Close()
		})
	}
}
