package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage/storage"
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
	storage storage.Storage
}

func NewHandler(storage storage.Storage) MetricHandler {
	InitMetricValidator()
	return MetricHandler{
		storage: storage,
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
	} else {
		h.storage.CounterUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *MetricHandler) Gauge(w http.ResponseWriter, r *http.Request) {
	name, valueString := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	value, err := strconv.ParseFloat(valueString, 64)

	if err != nil {
		http.Error(w, notFloatError, http.StatusBadRequest)
	} else {
		h.storage.GaugeUpdate(name, value)
		w.WriteHeader(http.StatusOK)
	}
}

func (h MetricHandler) MetricValue(w http.ResponseWriter, r *http.Request) {
	metricType, name := chi.URLParam(r, "type"), chi.URLParam(r, "name")

	if metricType == metric.GaugeTypeName {
		metric, ok := h.storage.GetGaugeMetric(name)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
		}

		w.Write([]byte(fmt.Sprint(metric)))
		w.WriteHeader(http.StatusOK)
	} else if metricType == metric.CounterTypeName {
		metric, ok := h.storage.GetCounterMetric(name)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
		}
		w.Write([]byte(fmt.Sprint(metric)))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, fmt.Sprintf(metricNotFoundError, name), http.StatusNotFound)
	}
}

func (h MetricHandler) MetricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	renderData := map[string]interface{}{
		"gaugeMetrics":   h.storage.GetGauge(),
		"counterMetrics": h.storage.GetCounter(),
	}

	template.HTMLTemplate.Execute(w, renderData)
}

func (h *MetricHandler) JSONUpdate(w http.ResponseWriter, r *http.Request) {
	var req Metric

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	if err := ValidateStruct(req); err != nil {
		http.Error(w, fmt.Sprintf(JSONValidationError, err.Error()), http.StatusBadRequest)
		return
	}

	if req.MType == metric.GaugeTypeName {
		h.storage.GaugeUpdate(req.ID, *req.Value)
		return
	}

	h.storage.CounterUpdate(req.ID, *req.Delta)
}

func (h MetricHandler) GetJSONMetric(w http.ResponseWriter, r *http.Request) {
	var req MetricReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, JSONDecodeError, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf(JSONValidationError, err.Error()), http.StatusBadRequest)
		return
	}

	body := Metric{
		ID:    req.ID,
		MType: req.MType,
	}

	if req.MType == metric.GaugeTypeName {
		metric, ok := h.storage.GetGaugeMetric(req.ID)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, req.ID), http.StatusNotFound)
			return
		}

		value := float64(metric)
		body.Value = &value
	} else {
		metric, ok := h.storage.GetCounterMetric(req.ID)

		if !ok {
			http.Error(w, fmt.Sprintf(metricNotFoundError, req.ID), http.StatusNotFound)
			return
		}

		value := int64(metric)
		body.Delta = &value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonValue, _ := json.Marshal(body)
	w.Write(jsonValue)
}
