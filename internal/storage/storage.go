package storage

import (
	"github.com/Kuart/metric-collector/internal/metric"
)

type GaugeStorage map[string]metric.GaugeValue

type CounterStorage map[string]metric.CounterValue

const (
	MetricNotFound = iota
	GaugeMetric
	CounterMetric
)

var gauge = &GaugeStorage{}
var counter = &CounterStorage{}

func GetGaugeStorage() GaugeStorage {
	gaugeStorage := make(map[string]metric.GaugeValue)

	for key, val := range *gauge {
		gaugeStorage[key] = val
	}

	return gaugeStorage
}

func GetCounterStorage() CounterStorage {
	counterStorage := make(map[string]metric.CounterValue)

	for key, val := range *counter {
		counterStorage[key] = val
	}

	return counterStorage
}

func GaugeUpdate(name string, value float64) {
	(*gauge)[name] += metric.GaugeValue(value)
}

func CounterUpdate(name string, value int64) {
	(*counter)[name] += metric.CounterValue(value)
}

func GetMetric(name string) (metric.GaugeValue, metric.CounterValue, int) {
	gaugeMetric, ok := (*gauge)[name]

	if ok {
		return gaugeMetric, 0, GaugeMetric
	}

	counterMetric, ok := (*counter)[name]

	if ok {
		return gaugeMetric, counterMetric, CounterMetric
	}

	return gaugeMetric, counterMetric, MetricNotFound
}

func CreateGauge() *GaugeStorage {
	return &GaugeStorage{}
}
