package main

import (
	"fmt"
	env "github.com/Kuart/metric-collector/cmd"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/template"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	template.SetupMetricTemplate()
	handler.InitMetricValidator()
	RunServer()
}

func RunServer() {
	r := chi.NewRouter()
	handler.SetRoutes(r)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", env.Host, env.Port),
		Handler: r,
	}

	server.ListenAndServe()
}
