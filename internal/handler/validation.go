package handler

import (
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/go-playground/validator/v10"
)

type MetricReq struct {
	ID    string `json:"id" validate:"required"`
	MType string `json:"type" validate:"required,oneof='gauge' 'counter'"`
}

var validate = validator.New()

func InitMetricValidator() {
	validate.RegisterStructValidation(metricValidation, metric.Metric{})
}

func metricValidation(sl validator.StructLevel) {
	cur := sl.Current().Interface()
	curMetric := cur.(metric.Metric)

	if curMetric.MType == metric.GaugeTypeName && curMetric.Value == nil {
		sl.ReportError(cur, "Value", "Value", "required", "")
	}

	if curMetric.MType == metric.CounterTypeName && curMetric.Delta == nil {
		sl.ReportError(cur, "Delta", "Delta", "required", "")
	}
}

func ValidateStruct(stc metric.Metric) error {
	if err := validate.Struct(stc); err != nil {
		return err
	}

	return nil
}
