package handler

import (
	"github.com/Kuart/metric-collector/internal/encryption"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
	crypto   encryption.Encryption
}

type MetricReq struct {
	ID    string `json:"id" validate:"required"`
	MType string `json:"type" validate:"required,oneof='gauge' 'counter'"`
	Hash  string `json:"hash,omitempty"`
}

func NewValidator(c encryption.Encryption) Validator {
	v := Validator{
		validate: validator.New(),
		crypto:   c,
	}

	v.validate.RegisterStructValidation(v.metricValidation, metric.Metric{})

	if c.IsKey {
		v.validate.RegisterStructValidation(v.hashValidation, metric.Metric{})
	}

	return v
}

func (v Validator) metricValidation(sl validator.StructLevel) {
	cur := sl.Current().Interface()
	curMetric := cur.(metric.Metric)

	if curMetric.MType == metric.GaugeTypeName && curMetric.Value == nil {
		sl.ReportError(cur, "Value", "Value", "required", "")
		return
	}

	if curMetric.MType == metric.CounterTypeName && curMetric.Delta == nil {
		sl.ReportError(cur, "Delta", "Delta", "required", "")
	}
}

func (v Validator) hashValidation(sl validator.StructLevel) {
	cur := sl.Current().Interface()
	curMetric := cur.(metric.Metric)

	if curMetric.Hash == "" {
		sl.ReportError(cur, "Hash", "Hash", "required", "")
		return
	}

	if !v.crypto.Compare(curMetric.Hash, curMetric) {
		sl.ReportError(cur, "Hash", "Hash", "not eq", "")
	}
}

func (v Validator) ValidateStruct(stc metric.Metric) error {
	if err := v.validate.Struct(stc); err != nil {
		return err
	}

	return nil
}
