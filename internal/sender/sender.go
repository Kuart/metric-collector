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

		url := fmt.Sprintf("%s/%s/%s/%f", c.updatePath, body.MType, body.ID, floatValue)

		c.doRequest(url)
		c.doJsonRequest(body)
	}
}

func (c *Client) sendCounter(counter metric.Counter) {
	intValue := int64(counter.Value)

	body := handler.Metric{
		ID:    counter.Name,
		MType: metric.CounterTypeName,
		Delta: &intValue,
	}

	url := fmt.Sprintf("%s/%s/%s/%d", c.updatePath, body.MType, body.ID, intValue)
	c.doRequest(url)
	c.doJsonRequest(body)
}

func (c *Client) doJsonRequest(body handler.Metric) {
	jsonValue, err := json.Marshal(body)

	if err != nil {
		log.Printf("%s metric not sended, json marshal err: %s", body.ID, err)
		return
	}

	buff := bytes.NewBuffer(jsonValue)
	response, err := c.client.Post(c.updatePath, "application/json;charset=utf-8", buff)

	if err != nil {
		log.Printf("%s metric not sended, err: %s", body.ID, err)
	}

	defer response.Body.Close()
}

func (c *Client) doRequest(url string) {
	response, err := c.client.Post(url, "text/plain", nil)

	if err != nil {
		log.Printf("%s metric not sended, err: %s", url, err)
	}

	defer response.Body.Close()
}
