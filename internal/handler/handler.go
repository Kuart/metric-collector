package handler

import (
	"fmt"
	"github.com/Kuart/metric-collector/internal/storage"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	notIntError    = "value is not int type"
	notFloatError  = "value is not float type"
	metricNotFound = "metric \"%s\" is not found"
)

func SetRoutes(r *chi.Mux) {
	r.Get("/value/{name}", MetricValueHandler)
	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{name}/{value}", CounterHandler)
		r.Post("/gauge/{name}/{value}", GaugeHandler)
	})
	r.Get("/", MetricsPageHandler)
}

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	name, valueString := getUrlParam(r)
	value, err := strconv.ParseInt(valueString, 10, 64)

	if err != nil {
		http.Error(w, notIntError, http.StatusBadRequest)
	} else {
		storage.CounterUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
	name, valueString := getUrlParam(r)
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		http.Error(w, notFloatError, http.StatusBadRequest)
	} else {
		storage.GaugeUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func MetricValueHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := getUrlParam(r)
	gauge, counter, status := storage.GetMetric(name)

	if status == storage.CounterMetric {
		w.Write([]byte(fmt.Sprint(counter)))
		w.WriteHeader(http.StatusOK)
	} else if status == storage.GaugeMetric {
		w.Write([]byte(fmt.Sprint(gauge)))
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

	template.HtmlTemplate.Execute(w, renderData)
}

func getUrlParam(r *http.Request) (string, string) {
	return chi.URLParam(r, "name"), chi.URLParam(r, "value")
}
