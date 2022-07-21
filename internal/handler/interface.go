package handler

import (
	"context"
	"github.com/Kuart/metric-collector/internal/metric"
)

type MetricHandler struct {
	controller Controller
	validator  Validator
	crypto     Encryption
}

type Controller interface {
	LoadToStorage()
	UpdateStorage(ctx context.Context, m metric.Metric) error
	GetMetric(ctx context.Context, m metric.Metric) (metric.Metric, bool)
	GetAllMetrics(ctx context.Context) (map[string]interface{}, error)
	GroupUpdateStorage(ctx context.Context, metrics []metric.Metric) error
	SaveToFile(ctx context.Context)
	PingDB() bool
	Close()
}

type Encryption interface {
	EncodeMetric(m metric.Metric) metric.Metric
	Compare(h string, m metric.Metric) bool
}
