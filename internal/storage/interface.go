package storage

import (
	"context"
	config "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage/inmemory"
)

type Controller struct {
	inmemory InmemoryStorage
	file     FileStorage
	db       DataBase
	sCfg     config.Config
	isSync   bool
	isUseDB  bool
}

type InmemoryStorage interface {
	GaugeUpdate(name string, value float64)
	CounterUpdate(name string, value int64)
	GetGauge() metric.GaugeState
	GetCounter() metric.CounterState
	GetGaugeMetric(name string) (metric.Metric, bool)
	GetCounterMetric(name string) (metric.Metric, bool)
	UpdateFromFile(data inmemory.FileStorage)
}

type FileStorage interface {
	GetFileData() (ifs inmemory.FileStorage, err error)
	Save(metrics map[string]interface{}) (err error)
	CloseFile()
}

type DataBase interface {
	Ping() bool
	Update(ctx context.Context, m metric.Metric) error
	GetMetric(ctx context.Context, m metric.Metric) (metric.Metric, bool)
	GetAllMetrics(ctx context.Context, MType string) (map[string]interface{}, error)
	BatchUpdate(metrics []metric.Metric) error
	Close()
}
