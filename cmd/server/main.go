package main

import (
	serverConfig "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/storage"
	"github.com/Kuart/metric-collector/internal/storage/file"
	"github.com/Kuart/metric-collector/internal/storage/inmemory"
	"github.com/Kuart/metric-collector/internal/template"
	"net/http"
)

func main() {
	config := serverConfig.New()
	inmemoryStorage := inmemory.New()
	fileStorage := file.New(config)
	srgController := storage.New(config, inmemoryStorage, fileStorage)
	template.SetupMetricTemplate()

	metricHandler := handler.NewHandler(srgController)
	r := handler.NewRouter(metricHandler)

	server := &http.Server{
		Addr:    config.Address,
		Handler: r,
	}

	server.ListenAndServe()

	defer fileStorage.CloseFile()
}
