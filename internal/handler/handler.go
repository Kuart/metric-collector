package handler

import (
	"fmt"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	notIntError     = "value is not int type"
	notFloatError   = "value is not float type"
	metricTypeError = "metric type not implemented"
	metricNotFound  = "metric \"%s\" is not found"
)

func SetRoutes(r *chi.Mux) {
	r.Get("/value/{type}/{name}", MetricValueHandler)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", UpdateHandler)
	})
	r.Get("/", MetricsPageHandler)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")

	switch metricType {
	case metric.GaugeTypeName:
		GaugeHandler(w, r)
	case metric.CounterTypeName:
		CounterHandler(w, r)
	default:
		http.Error(w, metricTypeError, http.StatusNotImplemented)
	}
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	name, valueString := chi.URLParam(r, "name"), chi.URLParam(r, "value")

	value, err := strconv.ParseInt(valueString, 10, 64)

	if err != nil {
		http.Error(w, notIntError, http.StatusBadRequest)
	} else {
		storage.CounterUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
	name, valueString := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		http.Error(w, notFloatError, http.StatusBadRequest)
	} else {
		storage.GaugeUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func MetricValueHandler(w http.ResponseWriter, r *http.Request) {
	metricType, name := chi.URLParam(r, "type"), chi.URLParam(r, "name")

	if metricType == metric.GaugeTypeName {
		metric, ok := storage.GetGaugeMetric(name)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFound, name), http.StatusNotFound)
		}

		w.Write([]byte(fmt.Sprint(metric)))
		w.WriteHeader(http.StatusOK)
	} else if metricType == metric.CounterTypeName {
		metric, ok := storage.GetCounterMetric(name)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFound, name), http.StatusNotFound)
		}
		w.Write([]byte(fmt.Sprint(metric)))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf(metricNotFound, name), http.StatusNotFound)
	}

}

func MetricsPageHandler(w http.ResponseWriter, r *http.Request) {
	renderData := map[string]interface{}{
		"gaugeMetrics":   storage.GetGaugeStorage(),
		"counterMetrics": storage.GetCounterStorage(),
	}

	template.HTMLTemplate.Execute(w, renderData)
}
