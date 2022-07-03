package metric

import (
	"math/rand"
	"runtime"
)

const (
	GaugeTypeName   = "gauge"
	CounterTypeName = "counter"
)

type Counter struct {
	Name  string
	Value int64
}

type Gauge struct {
	Name  string
	Value float64
}

type GaugeState = map[string]float64
type CounterState = map[string]int64

type Metric struct {
	ID    string   `json:"id" validate:"required"`
	MType string   `json:"type" validate:"required,oneof='gauge' 'counter'"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func NewMetricToSend[V float64 | int64](ID string, MType string, Value V) Metric {
	m := Metric{
		ID:    ID,
		MType: MType,
	}

	if MType == GaugeTypeName {
		v := float64(Value)
		m.Value = &v
	} else if MType == CounterTypeName {
		v := int64(Value)
		m.Delta = &v
	}

	return m
}

func GetGauge() []Gauge {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gauges := []Gauge{
		{"Alloc", float64(memStats.Alloc)},
		{"BuckHashSys", float64(memStats.BuckHashSys)},
		{"Frees", float64(memStats.Frees)},
		{"GCCPUFraction", float64(memStats.GCCPUFraction)},
		{"GCSys", float64(memStats.GCSys)},
		{"HeapAlloc", float64(memStats.HeapAlloc)},
		{"HeapIdle", float64(memStats.HeapIdle)},
		{"HeapInuse", float64(memStats.HeapInuse)},
		{"HeapObjects", float64(memStats.HeapObjects)},
		{"HeapReleased", float64(memStats.HeapReleased)},
		{"HeapSys", float64(memStats.HeapSys)},
		{"LastGC", float64(memStats.LastGC)},
		{"Lookups", float64(memStats.Lookups)},
		{"MCacheInuse", float64(memStats.MCacheInuse)},
		{"MCacheSys", float64(memStats.MCacheSys)},
		{"MSpanInuse", float64(memStats.MSpanInuse)},
		{"MSpanSys", float64(memStats.MSpanSys)},
		{"Mallocs", float64(memStats.Mallocs)},
		{"NextGC", float64(memStats.NextGC)},
		{"NumForcedGC", float64(memStats.NumForcedGC)},
		{"NumGC", float64(memStats.NumGC)},
		{"OtherSys", float64(memStats.OtherSys)},
		{"PauseTotalNs", float64(memStats.PauseTotalNs)},
		{"StackInuse", float64(memStats.StackInuse)},
		{"StackSys", float64(memStats.StackSys)},
		{"Sys", float64(memStats.Sys)},
		{"TotalAlloc", float64(memStats.TotalAlloc)},
	}

	return gauges
}

func GetRandomGauge() Gauge {
	return Gauge{"RandomValue", rand.Float64()}
}

func GetCounter(count int64) Counter {
	return Counter{"PollCount", count}
}

func (counter *Counter) PollTick() {
	counter.Value += 1
}

func (counter *Counter) Clear() {
	counter.Value = 0
}
