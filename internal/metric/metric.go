package metric

import (
	"math/rand"
	"runtime"
)

const (
	GaugeTypeName   = "gauge"
	CounterTypeName = "counter"
)

type CounterValue int64

type GaugeValue float64

type Counter struct {
	Name  string
	Value CounterValue
}

type Gauge struct {
	Name  string
	Value GaugeValue
}

func GetGauge() []Gauge {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	gauges := []Gauge{
		{"Alloc", GaugeValue(memStats.Alloc)},
		{"BuckHashSys", GaugeValue(memStats.BuckHashSys)},
		{"Frees", GaugeValue(memStats.Frees)},
		{"GCCPUFraction", GaugeValue(memStats.GCCPUFraction)},
		{"GCSys", GaugeValue(memStats.GCSys)},
		{"HeapAlloc", GaugeValue(memStats.HeapAlloc)},
		{"HeapIdle", GaugeValue(memStats.HeapIdle)},
		{"HeapInuse", GaugeValue(memStats.HeapInuse)},
		{"HeapObjects", GaugeValue(memStats.HeapObjects)},
		{"HeapReleased", GaugeValue(memStats.HeapReleased)},
		{"HeapSys", GaugeValue(memStats.HeapSys)},
		{"LastGC", GaugeValue(memStats.LastGC)},
		{"Lookups", GaugeValue(memStats.Lookups)},
		{"MCacheInuse", GaugeValue(memStats.MCacheInuse)},
		{"MCacheSys", GaugeValue(memStats.MCacheSys)},
		{"MSpanInuse", GaugeValue(memStats.MSpanInuse)},
		{"MSpanSys", GaugeValue(memStats.MSpanSys)},
		{"Mallocs", GaugeValue(memStats.Mallocs)},
		{"NextGC", GaugeValue(memStats.NextGC)},
		{"NumForcedGC", GaugeValue(memStats.NumForcedGC)},
		{"NumGC", GaugeValue(memStats.NumGC)},
		{"OtherSys", GaugeValue(memStats.OtherSys)},
		{"PauseTotalNs", GaugeValue(memStats.PauseTotalNs)},
		{"StackInuse", GaugeValue(memStats.StackInuse)},
		{"StackSys", GaugeValue(memStats.StackSys)},
		{"Sys", GaugeValue(memStats.Sys)},
		{"TotalAlloc", GaugeValue(memStats.TotalAlloc)},
	}

	return gauges
}

func GetRandomGauge() Gauge {
	return Gauge{"RandomValue", GaugeValue(rand.Float64())}
}

func GetCounter(count CounterValue) Counter {
	return Counter{"PollCount", count}
}

func (counter *Counter) PollTick() {
	counter.Value += 1
}

func (counter *Counter) Clear() {
	counter.Value = 0
}
