package storage

import (
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

	if sCfg.Restore {
		controller.LoadToStorage()
	}

	if sCfg.StoreInterval > 0 && sCfg.Address != "" {
		go controller.initFileSave(sCfg.StoreInterval)
	}

	if sCfg.StoreInterval == 0 {
		controller.isSync = true
	}

	if sCfg.DatabaseDSN != "" {
		dtb, err := database.New(sCfg)

		if err == nil {
			controller.db = dtb
			controller.isUseDB = true
		}
	}

	return controller
}

func (c Controller) initFileSave(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C
		c.SaveToFile()
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

func (c Controller) UpdateStorage(m metric.Metric) {
	if m.MType == metric.GaugeTypeName {
		c.inmemory.GaugeUpdate(m.ID, *m.Value)
	} else if m.MType == metric.CounterTypeName {
		c.inmemory.CounterUpdate(m.ID, *m.Delta)
	}

	if c.isSync {
		c.SaveToFile()
	}
}

func (c Controller) GetMetric(m metric.Metric) (metric.Metric, bool) {
	if m.MType == metric.GaugeTypeName {
		return c.inmemory.GetGaugeMetric(m.ID)
	}

	return c.inmemory.GetCounterMetric(m.ID)
}

func (c Controller) GetAllMetrics() map[string]interface{} {
	metrics := map[string]interface{}{
		"Gauge":   c.inmemory.GetGauge(),
		"Counter": c.inmemory.GetCounter(),
	}

	return metrics
}

func (c Controller) SaveToFile() {
	metrics := c.GetAllMetrics()
	c.file.Save(metrics)
}

func (c Controller) PingDB() bool {
	return c.db.Ping()
}

func (c Controller) Close() {
	c.file.CloseFile()
}
