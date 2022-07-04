package storage

import (
	"context"
	config "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage/database"
	"github.com/Kuart/metric-collector/internal/storage/file"
	"github.com/Kuart/metric-collector/internal/storage/inmemory"
	"html/template"
	"log"
	"time"
)

type Controller struct {
	inmemory inmemory.Storage
	file     file.Storage
	db       database.DB
	sCfg     config.Config
	isSync   bool
	isUseDB  bool
}

var HTMLTemplate *template.Template

func New(sCfg config.Config) Controller {
	controller := Controller{
		inmemory: inmemory.New(),
		file:     file.New(sCfg),
		sCfg:     sCfg,
	}

	if sCfg.DatabaseDSN != "" {
		dtb, err := database.New(sCfg)

		if err == nil {
			controller.db = dtb
			controller.isUseDB = true
			return controller
		}
	}

	if sCfg.Restore {
		controller.LoadToStorage()
	}

	if sCfg.StoreInterval > 0 && sCfg.Address != "" {
		go controller.initFileSave(sCfg.StoreInterval)
	}

	if sCfg.StoreInterval == 0 {
		controller.isSync = true
	}

	return controller
}

func (c Controller) initFileSave(interval time.Duration) {
	ctx := context.Background()
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C
		c.SaveToFile(ctx)
	}
}

func (c Controller) LoadToStorage() {
	data, err := c.file.GetFileData()

	if err != nil {
		log.Print(err)
		return
	}

	c.inmemory.UpdateFromFile(data)
}

func (c Controller) UpdateStorage(ctx context.Context, m metric.Metric) error {
	if c.isUseDB {
		return c.db.Update(ctx, m)
	}

	if m.MType == metric.GaugeTypeName {
		c.inmemory.GaugeUpdate(m.ID, *m.Value)
	} else if m.MType == metric.CounterTypeName {
		c.inmemory.CounterUpdate(m.ID, *m.Delta)
	}

	if c.isSync {
		c.SaveToFile(ctx)
	}

	return nil
}

func (c Controller) GetMetric(ctx context.Context, m metric.Metric) (metric.Metric, bool) {
	if c.isUseDB {
		return c.db.GetMetric(ctx, m)
	}

	if m.MType == metric.GaugeTypeName {
		return c.inmemory.GetGaugeMetric(m.ID)
	}

	return c.inmemory.GetCounterMetric(m.ID)
}

func (c Controller) GetAllMetrics(ctx context.Context) (map[string]interface{}, error) {
	if c.isUseDB {
		gauge, err := c.db.GetAllMetrics(ctx, metric.GaugeTypeName)

		if err != nil {
			log.Printf("get all gauge err: %s", err)
			return nil, err
		}

		counter, err := c.db.GetAllMetrics(ctx, metric.CounterTypeName)

		if err != nil {
			log.Printf("get all counter err: %s", err)
			return nil, err
		}

		metrics := map[string]interface{}{
			"Gauge":   gauge,
			"Counter": counter,
		}

		return metrics, nil
	}

	metrics := map[string]interface{}{
		"Gauge":   c.inmemory.GetGauge(),
		"Counter": c.inmemory.GetCounter(),
	}

	return metrics, nil
}

func (c Controller) SaveToFile(ctx context.Context) {
	metrics, err := c.GetAllMetrics(ctx)

	if err != nil {
		log.Printf("get all metrics err: %s", err)
	}

	c.file.Save(metrics)
}

func (c Controller) PingDB() bool {
	return c.db.Ping()
}

func (c Controller) Close() {
	c.file.CloseFile()
	c.db.Close()
}
