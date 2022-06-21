package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kuart/metric-collector/internal/handler"
	"github.com/Kuart/metric-collector/internal/metric"
	"log"
	"net/http"
	"time"
)

func NewMetricClient(address string, pollInterval time.Duration) *Client {
	updatePath := fmt.Sprintf("http://%s/update/", address)
	return &Client{
		updatePath: updatePath,
		client: &http.Client{
			Timeout: pollInterval,
		},
	}
}

func (c *Client) SendMetrics(gauge map[string]metric.GaugeValue, counter metric.Counter) {
	c.sendGauge(gauge)
	c.sendCounter(counter)
}

func (c *Client) sendGauge(gauge map[string]metric.GaugeValue) {
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
	jsonValue, err := json.Marshal(body)

	if err != nil {
		log.Printf("%s metric not sended, json marshal err: %s", body.ID, err)
		return
	}

	buff := bytes.NewBuffer(jsonValue)
	response, err := c.client.Post(c.updatePath, "application/json;charset=utf-8", buff)

	if err != nil {
		log.Printf("%s metric not sended, err: %s", body.ID, err)
		return
	}

	defer response.Body.Close()
}
