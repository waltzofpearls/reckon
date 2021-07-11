package prom

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/waltzofpearls/reckon/config"
	"go.uber.org/zap"
)

type Client struct {
	prom   v1.API
	config *config.Config
	logger *zap.Logger
}

func NewClient(cf *config.Config, lg *zap.Logger) *Client {
	return &Client{
		config: cf,
		logger: lg,
	}
}

func (c *Client) getOrInit() (v1.API, error) {
	if c.prom != nil {
		return c.prom, nil
	}

	if c.config.PromClientURL == "" {
		return nil, errors.New("PROM_CLIENT_URL cannot be empty")
	}

	tlsConfig, err := c.config.PromClientTLS()
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(api.Config{
		Address: c.config.PromClientURL,
		RoundTripper: &http.Transport{
			TLSClientConfig:   tlsConfig,
			DisableKeepAlives: true,
		},
	})
	if err != nil {
		return nil, err
	}

	c.prom = v1.NewAPI(client)
	return c.prom, nil
}

func (c *Client) GetCurrentMetricValue(ctx context.Context, query string) ([]Metric, error) {
	prom, err := c.getOrInit()
	if err != nil {
		return nil, err
	}
	value, warnings, err := prom.Query(ctx, query, time.Time{})
	if len(warnings) > 0 {
		c.logger.Warn("warnings returned from prometheus range query", zap.Strings("warnings", warnings))
	}
	if err != nil {
		c.logger.Error("failed querying prometheus metrics range", zap.String("query", query), zap.Error(err))
		return nil, errors.New("failed querying prometheus metrics range")
	}
	return c.metricsFrom(value)
}

func (c *Client) metricsFrom(value model.Value) ([]Metric, error) {
	vector := value.(model.Vector)
	if vector.Len() == 0 {
		c.logger.Warn("empty vector")
		return nil, nil
	}
	var metrics []Metric
	for _, sample := range vector {
		labels := make(map[string]string)
		for k, v := range sample.Metric {
			labels[string(k)] = string(v)
		}
		delete(labels, model.MetricNameLabel)
		metrics = append(metrics, Metric{
			Name:   string(sample.Metric[model.MetricNameLabel]),
			Labels: labels,
			Values: []SamplePair{
				{
					Value: float64(sample.Value),
					Time:  sample.Timestamp.Time(),
				},
			},
		})
	}
	return metrics, nil
}

func (c *Client) GetMetricRangeData(ctx context.Context, query string, start, end time.Time, chunkSize time.Duration) ([]Metric, error) {
	prom, err := c.getOrInit()
	if err != nil {
		return nil, err
	}
	if chunkSize == 0 {
		chunkSize = end.Sub(start)
	}
	start = start.Round(time.Second)
	end = end.Round(time.Second)
	chunkSize = chunkSize.Round(time.Second)
	if end.Sub(start).Seconds() < chunkSize.Seconds() {
		return nil, errors.New("specified chunkSize is too big")
	}

	metricMap := make(map[string]Metric)

	for start.Before(end) {
		to := start.Add(chunkSize)
		if to.After(end) {
			chunkSize = end.Sub(start)
		}
		chunked := fmt.Sprintf("%s[%vs]", query, chunkSize.Seconds())
		value, warnings, err := prom.Query(ctx, chunked, to)
		if len(warnings) > 0 {
			c.logger.Warn("warnings returned from prometheus range query", zap.Strings("warnings", warnings))
		}
		if err != nil {
			c.logger.Error("failed querying prometheus metrics range",
				zap.Time("start", start), zap.Time("end", end), zap.Duration("chunk_size", chunkSize),
				zap.String("query", chunked), zap.Time("to", to), zap.Error(err))
			return nil, errors.New("failed querying prometheus metrics range")
		}
		c.updateMetrics(value, metricMap)
		start = to
	}

	var metrics []Metric
	for _, metric := range metricMap {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (c *Client) updateMetrics(value model.Value, metrics map[string]Metric) {
	if value == nil {
		c.logger.Warn("model.Value is nil")
		return
	}
	matrix := value.(model.Matrix)
	if matrix.Len() == 0 {
		c.logger.Warn("empty matrix")
		return
	}
	for _, sample := range matrix {
		labels := make(map[string]string)
		for k, v := range sample.Metric {
			labels[string(k)] = string(v)
		}
		delete(labels, model.MetricNameLabel)
		metric := metrics[sample.Metric.String()]
		metric.Name = string(sample.Metric[model.MetricNameLabel])
		metric.Labels = labels
		for _, v := range sample.Values {
			metric.Values = append(metric.Values, SamplePair{
				Value: float64(v.Value),
				Time:  v.Timestamp.Time(),
			})
		}
		metrics[sample.Metric.String()] = metric
	}
}
