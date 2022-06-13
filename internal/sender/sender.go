package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/metric"
	"github.com/Kuart/metric-collector/internal/storage"
	"log"
	"net/http"
	"time"
)

func NewMetricClient(host string, port string, pollInterval time.Duration) *Client {
	basePath := fmt.Sprintf("http://%s:%s", host, port)
	updatePath := fmt.Sprintf("%s/update/", basePath)
	return &Client{
		updatePath: updatePath,
		client: &http.Client{
			Timeout: pollInterval,
		},
	}
}

func (c *Client) SendMetrics(gauge storage.GaugeStorage, counter metric.Counter) {
	c.sendGauge(gauge)
	c.sendCounter(counter)
}

func (c *Client) sendGauge(gauge storage.GaugeStorage) {
	for key, value := range gauge {
		floatValue := float64(value)

		body := handler.Metric{
			ID:    key,
			MType: metric.GaugeTypeName,
			Value: &floatValue,
		}

		c.doRequest(body)
	}
}

func (c *Client) sendCounter(counter metric.Counter) {
	intValue := int64(counter.Value)

	body := handler.Metric{
		ID:    counter.Name,
		MType: metric.CounterTypeName,
		Delta: &intValue,
	}

	c.doRequest(body)
}

func (c *Client) doRequest(body handler.Metric) {
	jsonValue, _ := json.Marshal(body)
	buff := bytes.NewBuffer(jsonValue)
	response, err := c.client.Post(c.updatePath, "application/json;charset=utf-8", buff)

	if err != nil {
		log.Printf("%s metric not sended, err: %s", body.ID, err)
	}

	defer response.Body.Close()
}
