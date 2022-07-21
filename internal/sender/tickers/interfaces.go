package tickers

import "github.com/Kuart/metric-collector/internal/metric"

type HTTPClient interface {
	SendMetrics(gauge metric.GaugeState, counter metric.Counter)
}

type TickerBuffer interface {
	WriteCommon()
	WriteGopsutil()
	Clear()
	Get() (metric.GaugeState, metric.Counter)
}
