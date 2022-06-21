package storage

import "github.com/Kuart/metric-collector/internal/metric"

type GaugeState = map[string]metric.GaugeValue
type CounterState = map[string]metric.CounterValue

type Storage struct {
	gauge   GaugeState
	counter CounterState
}

type FileStorage struct {
	Gauge   GaugeState
	Counter CounterState
}

func New() Storage {
	return Storage{
		gauge:   make(GaugeState),
		counter: make(CounterState),
	}
}

func (s *Storage) GaugeUpdate(name string, value float64) {
	s.gauge[name] = metric.GaugeValue(value)
}

func (s *Storage) CounterUpdate(name string, value int64) {
	s.counter[name] += metric.CounterValue(value)
}

func (s Storage) GetGauge() GaugeState {
	return s.gauge
}

func (s Storage) GetCounter() CounterState {
	return s.counter
}

func (s Storage) GetGaugeMetric(name string) (metric.GaugeValue, bool) {
	metric, ok := s.gauge[name]

	return metric, ok
}

func (s Storage) GetCounterMetric(name string) (metric.CounterValue, bool) {
	metric, ok := s.counter[name]

	return metric, ok
}

func (s *Storage) UpdateFromFile(data FileStorage) {
	for key, val := range data.Gauge {
		s.gauge[key] = val
	}

	for key, val := range data.Counter {
		s.counter[key] = val
	}
}
