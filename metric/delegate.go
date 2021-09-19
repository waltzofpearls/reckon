package metric

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/DataDog/go-python3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/model"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type delegate struct {
	logger *zap.Logger
	config *config.Config
	client *prom.Client

	original prom.Metric
	labels   []string
	key      string

	sem *semaphore.Weighted

	descs     map[string]*prometheus.Desc
	models    map[string]model.Trainer
	forecasts map[string]model.Forecasts

	modelNames      []string
	runtimeRegistry runtimeRegistry
}

func newDelegate(lg *zap.Logger, cf *config.Config, cl *prom.Client, data prom.Metric, modelNames []string) *delegate {
	labels := data.LabelNames()
	descs := make(map[string]*prometheus.Desc)
	models := make(map[string]model.Trainer)
	for _, modelName := range modelNames {
		mod, err := model.New(lg, modelName)
		if err != nil {
			lg.Error("cannot create model", zap.String("model", modelName),
				zap.String("original_metric", data.Name), zap.Error(err))
			continue
		}
		descs[modelName] = prometheus.NewDesc(
			data.Name+"_"+strings.ToLower(modelName),
			modelName+" forecasted metric value",
			append(labels, "column"),
			nil,
		)
		models[modelName] = mod
	}
	registry := newRuntimeRegistry(data.Name, modelNames)
	for key, metric := range registry {
		descs[key] = prometheus.NewDesc(
			metric.metric,
			metric.help,
			append(labels, "original_metric", "reckon_model", "reckon_host"),
			nil,
		)
	}
	return &delegate{
		logger: lg,
		config: cf,
		client: cl,

		original: data,
		labels:   labels,
		key:      data.Key(),

		sem: semaphore.NewWeighted(1),

		descs:     descs,
		models:    models,
		forecasts: make(map[string]model.Forecasts),

		modelNames: modelNames,
		// runtime metrics from each delegate
		// record number of model trainings, errors, durations and timestamps
		runtimeRegistry: registry,
	}
}

func (d *delegate) labelValues() []string {
	return d.original.LabelValuesFor(d.labels)
}

func (d *delegate) valuesFrom(modelName string, nearest time.Time) map[model.Column]float64 {
	values := make(map[model.Column]float64)
	if forecasts, exists := d.forecasts[modelName]; exists && len(forecasts) > 0 {
		forecast := forecasts.Nearest(nearest.Unix())
		if !forecast.IsZero() {
			for column, value := range forecast.Values {
				values[model.Column(column)] = value
			}
		}
	}
	return values
}

func (d *delegate) query() string {
	if len(d.original.Labels) == 0 {
		return d.original.Name
	}
	var labels []string
	for label, value := range d.original.Labels {
		labels = append(labels, fmt.Sprintf(`%s="%s"`, label, value))
	}
	return fmt.Sprintf("%s{%s}", d.original.Name, strings.Join(labels, ","))
}

func (d *delegate) train(ctx context.Context, module *python3.PyObject) {
	logger := d.logger.With(zap.String("prom_url", d.config.PromClientURL),
		zap.String("metric_name", d.original.Name), zap.Any("metric_labels", d.original.Labels))

	if !d.sem.TryAcquire(1) {
		logger.Info("skip training - metric delegate is already busy", zap.String("metric", d.key))
		return
	}
	defer d.sem.Release(1)

	duration := d.config.ChunkSize()
	if len(d.forecasts) == 0 { // initial run
		duration = d.config.RollingWindow
	}
	query := d.query()
	from := time.Now().UTC().Add(-duration)
	to := time.Now().UTC()

	fromField, toField := zap.Time("from", from), zap.Time("to", to)
	logger.Info("query prometheus and fetch metrics data", fromField, toField)
	d.runtimeRegistry.nowAll(d.modelNames, "reckon_prometheus_client_scrape_time_seconds")
	d.runtimeRegistry.incAll(d.modelNames, "reckon_prometheus_client_scrape_total")
	d.runtimeRegistry.setAll(d.modelNames, "reckon_data_scraped_duration_minutes", duration.Minutes())

	data, err := d.client.GetMetricRangeData(ctx, query, from, to, 0)
	if err != nil {
		logger.Error("failed querying prometheus data range", fromField, toField, zap.Error(err))
		d.runtimeRegistry.incAll(d.modelNames, "reckon_prometheus_client_scrape_errors_total")
		return
	}
	if len(data) == 0 {
		logger.Info("no data from range query", fromField, toField, zap.Error(err))
		return
	}
	d.runtimeRegistry.setAll(d.modelNames, "reckon_data_scraped_values", float64(len(data[0].Values)))
	if err := d.original.Append(data[0], d.config.RollingWindow); err != nil {
		logger.Error("failed appending new data to original metric data",
			zap.String("other_name", data[0].Name), zap.Any("other_labels", data[0].Labels), zap.Error(err))
		d.runtimeRegistry.incAll(d.modelNames, "reckon_prometheus_client_scrape_errors_total")
		return
	}
	d.runtimeRegistry.setAll(d.modelNames, "reckon_training_data_values", float64(len(d.original.Values)))

	var wg sync.WaitGroup
	for modelName, mod := range d.models {
		duration = config.RoundUpDuration(d.original.End().Sub(d.original.Start()), time.Minute)
		logger.Info("calling model train method", zap.String("model_name", modelName),
			zap.Int("length", len(d.original.Values)), zap.String("rolling_window", d.config.RollingWindow.String()),
			zap.String("data_range", duration.String()), zap.String("start", d.original.Start().String()),
			zap.String("end", d.original.End().String()))
		d.runtimeRegistry.now(modelName, "reckon_model_train_time_seconds")
		d.runtimeRegistry.inc(modelName, "reckon_model_train_total")
		d.runtimeRegistry.set(modelName, "reckon_training_data_duration_minutes", duration.Minutes())
		now := time.Now()
		wg.Add(1)
		go func(modelName string, mod model.Trainer) {
			defer wg.Done()
			defer func() {
				d.runtimeRegistry.set(modelName, "reckon_model_train_duration_seconds", time.Since(now).Seconds())
			}()
			d.forecasts[modelName], err = mod.Train(ctx, module, d.original, duration)
			if err != nil {
				logger.Error("unable to train model", zap.Error(err))
				d.runtimeRegistry.inc(modelName, "reckon_model_train_errors_total")
				return
			}
			logger.Info("received forecasted data", zap.Int("length", len(d.forecasts[modelName])))
			d.runtimeRegistry.set(modelName, "reckon_forecast_data_duration_minutes", duration.Minutes())
			d.runtimeRegistry.set(modelName, "reckon_forecast_data_values", float64(len(d.forecasts[modelName])))
		}(modelName, mod)
		d.runtimeRegistry.now(modelName, "reckon_forecast_data_received_time_seconds")
	}
	wg.Wait()
}
