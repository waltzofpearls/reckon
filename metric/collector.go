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
	build  *config.BuildInfo
	store  *Store
}

func NewCollector(cf *config.Config, lg *zap.Logger, bi *config.BuildInfo, st *Store) prometheus.Collector {
	return &Collector{
		config: cf,
		logger: lg,
		build:  bi,
		store:  st,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.logger.Info("describe reckon exporter metrics")
	c.store.ForEach(func(key string, del *delegate) {
		for _, desc := range del.descs {
			ch <- desc
		}
	})
	c.build.Describe(ch)
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Info("scrape reckon exporter metrics")
	c.store.ForEach(func(key string, del *delegate) {
		for modelOrRuntimeMetric, desc := range del.descs {
			if runtimeMetric, exists := del.runtimeRegistry[modelOrRuntimeMetric]; exists {
				// desc is a runtime metric
				// modelOrRuntimeMetric = {model_name}::{runtime_metric_name}
				ch <- prometheus.MustNewConstMetric(
					desc,
					runtimeMetric.typ,
					runtimeMetric.into(),
					append(del.labelValues(), runtimeMetric.labels()...)...,
				)
			} else {
				// desc is a forecast metric
				modelName := modelOrRuntimeMetric
				values := del.valuesFrom(modelName, time.Now().In(c.config.Location()))
				for column, value := range values {
					ch <- prometheus.MustNewConstMetric(
						desc,
						prometheus.GaugeValue,
						value,
						append(del.labelValues(), string(column))...,
					)
				}
			}
		}
		del.runtimeRegistry.nowAll(del.modelNames, "reckon_exporter_scraped_time_seconds")
		del.runtimeRegistry.incAll(del.modelNames, "reckon_exporter_scraped_total")
		c.store.Save(key, del)
	})
	c.build.Collect(ch)
}
