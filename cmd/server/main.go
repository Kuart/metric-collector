package main

import (
	"fmt"
	"github.com/Kuart/metric-collector/internal/api"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	template.SetupMetricTemplate()
	RunServer()
}

func RunServer() {
	r := chi.NewRouter()
	handler.SetRoutes(r)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", api.Host, api.Port),
		Handler: r,
	}

	server.ListenAndServe()
}
