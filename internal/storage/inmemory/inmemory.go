package inmemory

import "github.com/Kuart/metric-collector/internal/metric"

type Storage struct {
	gauge   metric.GaugeState
	counter metric.CounterState
}

type FileStorage struct {
	Gauge   metric.GaugeState
	Counter metric.CounterState
}

func New() Storage {
	return Storage{
		gauge:   make(metric.GaugeState),
		counter: make(metric.CounterState),
	}
}

func (s *Storage) GaugeUpdate(name string, value float64) {
	s.gauge[name] = value
}

func (s *Storage) CounterUpdate(name string, value int64) {
	s.counter[name] += value
}

func (s Storage) GetGauge() metric.GaugeState {
	return s.gauge
}

func (s Storage) GetCounter() metric.CounterState {
	return s.counter
}

func (s Storage) GetGaugeMetric(name string) (metric.Metric, bool) {
	val, ok := s.gauge[name]

	m := metric.Metric{
		ID:    name,
		Value: &val,
	}

	return m, ok
}

func (s Storage) GetCounterMetric(name string) (metric.Metric, bool) {
	val, ok := s.counter[name]

	m := metric.Metric{
		ID:    name,
		Delta: &val,
	}

	return m, ok
}

func (s *Storage) UpdateFromFile(data FileStorage) {
	for key, val := range data.Gauge {
		s.gauge[key] = val
	}

	for key, val := range data.Counter {
		s.counter[key] = val
	}
}
