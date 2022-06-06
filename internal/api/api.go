package api

import (
	"fmt"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage"
	"log"
	"net/http"
)

const (
	Host = "127.0.0.1"
	Port = "8080"
)

func SendMetric(client *http.Client, gauge storage.GaugeStorage, counter metric.Counter) {
	parseGauge(client, gauge)
	sendCounter(client, counter)
}

func parseGauge(client *http.Client, gauge storage.GaugeStorage) {
	for key, value := range gauge {
		sendGauge(client, key, value)
	}
}

func sendCounter(client *http.Client, counter metric.Counter) {
	pattern := "http://%s:%s/update/%s/%s/%d"
	url := fmt.Sprintf(pattern, Host, Port, metric.CounterTypeName, counter.Name, int64(counter.Value))
	response, err := client.Post(url, "text/plain", nil)

	if err != nil {
		log.Println(fmt.Sprintf("%s metric not sended, err: %s", counter.Name, err))
	}

	defer response.Body.Close()
}

func sendGauge(client *http.Client, name string, value metric.GaugeValue) {
	pattern := "http://%s:%s/update/%s/%s/%f"
	url := fmt.Sprintf(pattern, Host, Port, metric.GaugeTypeName, name, float64(value))
	response, err := client.Post(url, "text/plain", nil)

	if err != nil {
		log.Println(fmt.Sprintf("%s metric not sended, err: %s", name, err))
	}

	defer response.Body.Close()
}
