package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/encryption"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	notIntError         = "value is not int type"
	notFloatError       = "value is not float type"
	metricTypeError     = "metric type not implemented"
	metricNotFoundError = "metric \"%s\" is not found"
	metricsGetError     = "error while getting all metrics"
	JSONDecodeError     = "error in JSON decode"
	storageUpdateError  = "error during storage update"
	JSONValidationError = "JSON validation fail: \"%s\""
)

func NewHandler(controller Controller, c encryption.Encryption) MetricHandler {
	return MetricHandler{
		controller: controller,
		crypto:     c,
		validator:  NewValidator(c),
	}
}

func (h MetricHandler) Update(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")

	switch metricType {
	case metric.GaugeTypeName:
		h.Gauge(w, r)
	case metric.CounterTypeName:
		h.Counter(w, r)
	default:
		http.Error(w, metricTypeError, http.StatusNotImplemented)
	}
}

func (h *MetricHandler) Counter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	name, valueString := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	value, err := strconv.ParseInt(valueString, 10, 64)

	if err != nil {
		http.Error(w, notIntError, http.StatusBadRequest)
		return
	}

	m := metric.Metric{
		ID:    name,
		MType: metric.CounterTypeName,
		Delta: &value,
	}

	err = h.controller.UpdateStorage(ctx, m)

	if err != nil {
		log.Printf(storageUpdateError+":%s", err.Error())
		http.Error(w, storageUpdateError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) Gauge(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	name, valueString := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		http.Error(w, notFloatError, http.StatusBadRequest)
		return
	}

	m := metric.Metric{
		ID:    name,
		MType: metric.GaugeTypeName,
		Value: &value,
	}

	err = h.controller.UpdateStorage(ctx, m)

	if err != nil {
		log.Printf(storageUpdateError+":%s", err.Error())
		http.Error(w, storageUpdateError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (h MetricHandler) MetricValue(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	metricType, name := chi.URLParam(r, "type"), chi.URLParam(r, "name")

	if metricType == metric.GaugeTypeName {
		metric, ok := h.controller.GetMetric(ctx, metric.Metric{ID: name, MType: metric.GaugeTypeName})

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprint(*metric.Value)))
		w.WriteHeader(http.StatusOK)
	} else if metricType == metric.CounterTypeName {
		metric, ok := h.controller.GetMetric(ctx, metric.Metric{ID: name, MType: metric.CounterTypeName})

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprint(*metric.Delta)))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
	}
}

func (h MetricHandler) MetricsPage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	metrics, err := h.controller.GetAllMetrics(ctx)

	if err != nil {
		http.Error(w, fmt.Sprintf(metricsGetError+": %s", err), http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	template.HTMLTemplate.Execute(w, metrics)
}

func (h *MetricHandler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var req metric.Metric

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateStruct(req); err != nil {
		http.Error(w, fmt.Sprintf(JSONValidationError, err.Error()), http.StatusBadRequest)
		return
	}

	m := metric.Metric{ID: req.ID, MType: req.MType}

	if req.MType == metric.GaugeTypeName {
		m.Value = req.Value
		err := h.controller.UpdateStorage(ctx, m)

		if err != nil {
			log.Printf(storageUpdateError+":%s", err.Error())
			http.Error(w, storageUpdateError, http.StatusInternalServerError)
		}
		return
	}

	m.Delta = req.Delta
	err := h.controller.UpdateStorage(ctx, m)

	if err != nil {
		log.Printf(storageUpdateError+":%s", err.Error())
		http.Error(w, storageUpdateError, http.StatusInternalServerError)
	}
}

func (h MetricHandler) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var req MetricReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	if err := h.validator.validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf(JSONValidationError, err.Error()), http.StatusBadRequest)
		return
	}

	m, ok := h.controller.GetMetric(ctx, metric.Metric{
		ID:    req.ID,
		MType: req.MType,
	})

	if !ok {
		http.Error(w, fmt.Sprintf(metricNotFoundError, req.ID), http.StatusNotFound)
		return
	}

	result := h.crypto.EncodeMetric(m)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
	w.WriteHeader(http.StatusOK)
}

func (h MetricHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	ok := h.controller.PingDB()

	if !ok {
		http.Error(w, "database not responding", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h MetricHandler) GroupUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var metrics []metric.Metric
	var approvedMetrics []metric.Metric

	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	for _, mtr := range metrics {
		if err := h.validator.validate.Struct(mtr); err != nil {
			log.Printf("metric '%s' is not valid: %s", mtr.ID, err.Error())
			continue
		}
		approvedMetrics = append(approvedMetrics, mtr)
	}

	err := h.controller.GroupUpdateStorage(ctx, approvedMetrics)

	if err != nil {
		log.Printf("group update fail %s", err.Error())
		http.Error(w, "group update fail", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
