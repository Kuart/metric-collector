package main

import (
	"fmt"
	"github.com/Kuart/metric-collector/cmd"
	serverConfig "github.com/Kuart/metric-collector/config/server"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/storage/file"
	"github.com/Kuart/metric-collector/internal/storage/storage"
	"github.com/Kuart/metric-collector/internal/template"
	"net/http"
)

func main() {
	config := serverConfig.New()
	storage := storage.New()

	fileStorage := file.New(config.StoreFile, config.StoreInterval, &storage)
	fileStorage.LoadToStorage(config.Restore)
	go fileStorage.InitSaver()

	template.SetupMetricTemplate()

	metricHandler := handler.NewHandler(storage)
	r := handler.NewRouter(metricHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cmd.Host, cmd.Port),
		Handler: r,
	}

	server.ListenAndServe()
}
