package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler MetricHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Compress(5, "gzip"))

	r.Get("/value/{type}/{name}", handler.MetricValue)
	r.Post("/value", handler.GetJSONMetric)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handler.Update)
		r.Post("/", handler.JSONUpdate)
	})
	r.Get("/", handler.MetricsPage)

	return r
}
