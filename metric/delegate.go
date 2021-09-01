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

type Delegate struct {
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
}

func NewDelegate(lg *zap.Logger, cf *config.Config, cl *prom.Client, data prom.Metric) *Delegate {
	labels := data.LabelNames()
	descs := make(map[string]*prometheus.Desc)
	models := make(map[string]model.Trainer)
	for _, modelName := range cf.Models {
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
	return &Delegate{
		logger:    lg,
		config:    cf,
		client:    cl,
		original:  data,
		labels:    labels,
		key:       data.Key(),
		sem:       semaphore.NewWeighted(1),
		descs:     descs,
		models:    models,
		forecasts: make(map[string]model.Forecasts),
	}
}

func (d *Delegate) Key() string {
	return d.key
}

func (d *Delegate) Descs() map[string]*prometheus.Desc {
	return d.descs
}

func (d *Delegate) LabelValues() []string {
	return d.original.LabelValuesFor(d.labels)
}

func (d *Delegate) ValuesFrom(modelName string, nearest time.Time) map[model.Column]float64 {
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

func (d *Delegate) Query() string {
	if len(d.original.Labels) == 0 {
		return d.original.Name
	}
	var labels []string
	for label, value := range d.original.Labels {
		labels = append(labels, fmt.Sprintf(`%s="%s"`, label, value))
	}
	return fmt.Sprintf("%s{%s}", d.original.Name, strings.Join(labels, ","))
}

func (d *Delegate) Train(ctx context.Context, module *python3.PyObject) {
	logger := d.logger.With(zap.String("prom_url", d.config.PromClientURL),
		zap.String("metric_name", d.original.Name), zap.Any("metric_labels", d.original.Labels))

	if !d.sem.TryAcquire(1) {
		logger.Info("skip training - metric delegate is already busy", zap.String("metric", d.Key()))
		return
	}
	defer d.sem.Release(1)

	duration := d.config.ChunkSize()
	if len(d.forecasts) == 0 { // initial run
		duration = d.config.RollingWindow
	}
	query := d.Query()
	from := time.Now().UTC().Add(-duration)
	to := time.Now().UTC()

	fromField, toField := zap.Time("from", from), zap.Time("to", to)
	logger.Info("query prometheus and fetch metrics data", fromField, toField)

	data, err := d.client.GetMetricRangeData(ctx, query, from, to, 0)
	if err != nil {
		logger.Error("failed querying prometheus data range", fromField, toField, zap.Error(err))
		return
	}
	if len(data) == 0 {
		logger.Info("no data from range query", fromField, toField, zap.Error(err))
		return
	}
	if err := d.original.Append(data[0], d.config.RollingWindow); err != nil {
		logger.Error("failed appending new data to original metric data",
			zap.String("other_name", data[0].Name), zap.Any("other_labels", data[0].Labels), zap.Error(err))
		return
	}

	var wg sync.WaitGroup
	for modelName, mod := range d.models {
		duration = config.RoundUpDuration(d.original.End().Sub(d.original.Start()), time.Minute)
		logger.Info("calling model train method", zap.String("model_name", modelName),
			zap.Int("length", len(d.original.Values)), zap.String("rolling_window", d.config.RollingWindow.String()),
			zap.String("data_range", duration.String()), zap.String("start", d.original.Start().String()),
			zap.String("end", d.original.End().String()))
		wg.Add(1)
		go func(modelName string, mod model.Trainer) {
			defer wg.Done()
			d.forecasts[modelName] = mod.Train(ctx, module, d.original, duration)
			logger.Info("received forecasted data", zap.Int("length", len(d.forecasts[modelName])))
		}(modelName, mod)
	}
	wg.Wait()
}
