package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/encryption"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const (
	notIntError         = "value is not int type"
	notFloatError       = "value is not float type"
	metricTypeError     = "metric type not implemented"
	metricNotFoundError = "metric \"%s\" is not found"
	JSONDecodeError     = "error in JSON decode"
	JSONValidationError = "JSON validation fail: \"%s\""
)

type MetricHandler struct {
	controller storage.Controller
	validator  Validator
	crypto     encryption.Encryption
}

func NewHandler(controller storage.Controller, c encryption.Encryption) MetricHandler {
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

	h.controller.UpdateStorage(m)
	w.WriteHeader(http.StatusOK)
}

func (h *MetricHandler) Gauge(w http.ResponseWriter, r *http.Request) {
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

	h.controller.UpdateStorage(m)
	w.WriteHeader(http.StatusOK)

}

func (h MetricHandler) MetricValue(w http.ResponseWriter, r *http.Request) {
	metricType, name := chi.URLParam(r, "type"), chi.URLParam(r, "name")

	if metricType == metric.GaugeTypeName {
		metric, ok := h.controller.GetMetric(metric.Metric{ID: name, MType: metric.GaugeTypeName})

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
			return
		}

		w.Write([]byte(fmt.Sprint(*metric.Value)))
		w.WriteHeader(http.StatusOK)
	} else if metricType == metric.CounterTypeName {
		metric, ok := h.controller.GetMetric(metric.Metric{ID: name, MType: metric.CounterTypeName})

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
	w.Header().Set("Content-Type", "text/html")
	template.HTMLTemplate.Execute(w, h.controller.GetAllMetrics())
}

func (h *MetricHandler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
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
		h.controller.UpdateStorage(m)
		return
	}

	m.Delta = req.Delta
	h.controller.UpdateStorage(m)
}

func (h MetricHandler) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	var req MetricReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	if err := h.validator.validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf(JSONValidationError, err.Error()), http.StatusBadRequest)
		return
	}

	m, ok := h.controller.GetMetric(metric.Metric{
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
