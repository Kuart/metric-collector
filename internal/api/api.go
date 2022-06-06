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
	sendGauge(client, gauge)
	sendCounter(client, counter)
}

func sendGauge(client *http.Client, gauge storage.GaugeStorage) {
	for key, value := range gauge {
		send(client, metric.GaugeTypeName, key, value)
	}
}

func sendCounter(client *http.Client, counter metric.Counter) {
	send(client, metric.CounterTypeName, counter.Name, counter.Value)
}

func send[T metric.GaugeValue | metric.CounterValue](client *http.Client, metricType string, name string, value T) {
	url := createURL(metricType, name, value)
	response, err := client.Post(url, "text/plain", nil)

	if err != nil {
		log.Println(fmt.Sprintf("%s metric not sended, err: %s", name, err))
	}

	defer response.Body.Close()
}

func createURL[T metric.GaugeValue | metric.CounterValue](metricType string, name string, value T) string {
	if metricType == metric.GaugeTypeName {
		return fmt.Sprintf("http://%s:%s/update/%s/%s/%f", Host, Port, metricType, name, float64(value))
	}

	return fmt.Sprintf("http://%s:%s/update/%s/%s/%d", Host, Port, metricType, name, int64(value))
}
