package prom

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	promAPI "github.com/prometheus/client_golang/api"
	promAPIV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	reckonAPI "github.com/waltzofpearls/reckon/api"
	"github.com/waltzofpearls/reckon/config"
)

type Client struct {
	api promAPIV1.API

	*config.Config
}

func NewClient(config *config.Config) *Client {
	return &Client{
		Config: config,
	}
}

func (c *Client) getAPI() (promAPIV1.API, error) {
	if c.api != nil {
		return c.api, nil
	}

	tlsConfig, err := c.PromClientTLS()
	if err != nil {
		return nil, err
	}

	client, err := promAPI.NewClient(promAPI.Config{
		Address: c.PromClientURL,
		RoundTripper: &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return promAPIV1.NewAPI(client), nil
}

type labelSet map[string]string

func (c *Client) GetCurrentMetricValue(
	ctx context.Context, metricName string, labels labelSet,
) (reckonAPI.Metric, error) {
	clientAPI, err := c.getAPI()
	if err != nil {
		return reckonAPI.Metric{}, err
	}
	query := createQueryFrom(metricName, labels)
	value, _, _ := clientAPI.Query(ctx, query, time.Time{})
	return vectorToMetric(value)
}

func createQueryFrom(metricName string, labels labelSet) string {
	query := metricName
	if len(labels) > 0 {
		list := make([]string, len(labels))
		for k, v := range labels {
			list = append(list, fmt.Sprintf("%[1]s='%[2]s'", k, v))
		}
		query += fmt.Sprintf("{%s}", strings.Join(list, ", "))
	}
	return query
}

func vectorToMetric(value model.Value) (reckonAPI.Metric, error) {
	var metric reckonAPI.Metric
	vector := value.(model.Vector)
	if vector.Len() == 0 {
		return metric, nil
	}
	sample := vector[0]
	timestamp, err := ptypes.TimestampProto(sample.Timestamp.Time())
	if err != nil {
		return metric, err
	}
	labels := make(map[string]string)
	for k, v := range sample.Metric {
		labels[string(k)] = string(v)
	}
	delete(labels, model.MetricNameLabel)
	metric = reckonAPI.Metric{
		Name:   string(sample.Metric[model.MetricNameLabel]),
		Labels: labels,
		Values: []*reckonAPI.SamplePair{
			&reckonAPI.SamplePair{
				Value: float64(sample.Value),
				Time:  timestamp,
			},
		},
	}
	return metric, nil
}

func (c *Client) GetMetricRangeData(
	ctx context.Context,
	metricName string,
	labels labelSet,
	start, end time.Time,
	step time.Duration,
	chunkSize time.Duration,
) (reckonAPI.Metric, error) {
	var metric reckonAPI.Metric
	clientAPI, err := c.getAPI()
	if err != nil {
		return metric, err
	}
	if chunkSize == 0 {
		chunkSize = end.Sub(start)
	}
	start = start.Round(time.Second)
	end = end.Round(time.Second)
	chunkSize = chunkSize.Round(time.Second)
	if end.Sub(start).Seconds() < chunkSize.Seconds() {
		return metric, errors.New("Error: specified chunkSize is too big")
	}
	query := createQueryFrom(metricName, labels)
	for start.Before(end) {
		if start.Add(chunkSize).After(end) {
			chunkSize = end.Sub(start)
		}
		value, _, _ := clientAPI.QueryRange(ctx, query, promAPIV1.Range{
			Start: start,
			End:   end,
			Step:  step,
		})
		updateMetricWith(value, &metric)
		start = start.Add(chunkSize)
	}
	return metric, nil
}

func updateMetricWith(value model.Value, metric *reckonAPI.Metric) {
	matrix := value.(model.Matrix)
	if matrix.Len() == 0 {
		log.Println("empty matrix")
		return
	}
	sampleStream := matrix[0]
	labels := make(map[string]string)
	for k, v := range sampleStream.Metric {
		labels[string(k)] = string(v)
	}
	delete(labels, model.MetricNameLabel)
	metric.Name = string(sampleStream.Metric[model.MetricNameLabel])
	metric.Labels = labels
	for _, v := range sampleStream.Values {
		timestamp, err := ptypes.TimestampProto(v.Timestamp.Time())
		if err != nil {
			log.Println("failed to convert time to proto timestamp")
			continue
		}
		metric.Values = append(metric.Values, &reckonAPI.SamplePair{
			Value: float64(v.Value),
			Time:  timestamp,
		})
	}
}
