package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/waltzofpearls/reckon/config"
	"go.uber.org/zap"
)

type Collector struct {
	config *config.Config
	logger *zap.Logger
	store  *Store
}

func NewCollector(cf *config.Config, lg *zap.Logger, st *Store) prometheus.Collector {
	return &Collector{
		config: cf,
		logger: lg,
		store:  st,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.logger.Info("describe prometheus metrics")
	c.store.ForEach(func(key string, delegate *Delegate) {
		for _, desc := range delegate.Descs() {
			ch <- desc
		}
	})
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Info("collect prometheus metrics")
	c.store.ForEach(func(key string, delegate *Delegate) {
		for modelName, desc := range delegate.Descs() {
			values := delegate.ValuesFrom(modelName, time.Now().In(c.config.Location()))
			for column, value := range values {
				ch <- prometheus.MustNewConstMetric(
					desc,
					prometheus.GaugeValue,
					value,
					append(delegate.LabelValues(), string(column))...,
				)
			}
		}
	})
}
