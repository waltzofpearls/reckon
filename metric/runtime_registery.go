package metric

import (
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type runtimeMetric struct {
	metric   string
	help     string
	valBits  uint64
	typ      prometheus.ValueType
	original string
	model    string
}

func (r *runtimeMetric) labels() []string {
	hostname, _ := os.Hostname()
	// labels: original_metric, reckon_model, reckon_host
	return []string{r.original, r.model, hostname}
}

func (r *runtimeMetric) into() float64 {
	return math.Float64frombits(atomic.LoadUint64(&r.valBits))
}

func (r *runtimeMetric) set(value float64) {
	atomic.StoreUint64(&r.valBits, math.Float64bits(value))
}

func (r *runtimeMetric) add(val float64) {
	for {
		oldBits := atomic.LoadUint64(&r.valBits)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&r.valBits, oldBits, newBits) {
			return
		}
	}
}

func (r *runtimeMetric) inc() {
	r.add(1)
}

type runtimeRegistry map[string]*runtimeMetric

func newRuntimeRegistry(original string, models []string) runtimeRegistry {
	registry := make(runtimeRegistry)
	for _, modelName := range models {
		registry[modelName+"::reckon_prometheus_client_scrape_time_seconds"] = &runtimeMetric{
			metric:   "reckon_prometheus_client_scrape_time_seconds",
			help:     "timestamp of the last prometheus client scrape",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_exporter_scraped_time_seconds"] = &runtimeMetric{
			metric:   "reckon_exporter_scraped_time_seconds",
			help:     "timestamp of the last time reckon exporter scraped by prometheus",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_model_train_time_seconds"] = &runtimeMetric{
			metric:   "reckon_model_train_time_seconds",
			help:     "timestamp of the last reckon model training",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_forecast_data_received_time_seconds"] = &runtimeMetric{
			metric:   "reckon_forecast_data_received_time_seconds",
			help:     "timestamp of the last time receiving forecast data",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_model_train_total"] = &runtimeMetric{
			metric:   "reckon_model_train_total",
			help:     "number of times calling model train",
			typ:      prometheus.CounterValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_model_train_duration_seconds"] = &runtimeMetric{
			metric:   "reckon_model_train_duration_seconds",
			help:     "time taken in seconds from the last model training",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_model_train_errors_total"] = &runtimeMetric{
			metric:   "reckon_model_train_errors_total",
			help:     "number of model training errors",
			typ:      prometheus.CounterValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_prometheus_client_scrape_errors_total"] = &runtimeMetric{
			metric:   "reckon_prometheus_client_scrape_errors_total",
			help:     "prometheus client scraping errors",
			typ:      prometheus.CounterValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_exporter_scraped_total"] = &runtimeMetric{
			metric:   "reckon_exporter_scraped_total",
			help:     "number of times reckon exporter scraped by prometheus",
			typ:      prometheus.CounterValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_prometheus_client_scrape_total"] = &runtimeMetric{
			metric:   "reckon_prometheus_client_scrape_total",
			help:     "number of prometheus client scrape",
			typ:      prometheus.CounterValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_data_scraped_duration_minutes"] = &runtimeMetric{
			metric:   "reckon_data_scraped_duration_minutes",
			help:     "duration of data scraped from prometheus",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_forecast_data_duration_minutes"] = &runtimeMetric{
			metric:   "reckon_forecast_data_duration_minutes",
			help:     "duration of data being kept in memory",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_training_data_duration_minutes"] = &runtimeMetric{
			metric:   "reckon_training_data_duration_minutes",
			help:     "duration of data sent to model for training",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_data_scraped_values"] = &runtimeMetric{
			metric:   "reckon_data_scraped_values",
			help:     "number of the last scraped data points",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_forecast_data_values"] = &runtimeMetric{
			metric:   "reckon_forecast_data_values",
			help:     "number of the existing data points kept in memory",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
		registry[modelName+"::reckon_training_data_values"] = &runtimeMetric{
			metric:   "reckon_training_data_values",
			help:     "number of data points recently sent to model for training",
			typ:      prometheus.GaugeValue,
			original: original,
			model:    modelName,
		}
	}
	return registry
}

func (r runtimeRegistry) set(model, metric string, value float64) error {
	key := model + "::" + metric
	if _, exists := r[key]; exists {
		r[key].set(value)
		return nil
	}
	return fmt.Errorf("unable to find runtime metric %s", key)
}

func (r runtimeRegistry) setAll(models []string, metric string, value float64) error {
	var err error
	for _, modelName := range models {
		if err = r.set(modelName, metric, value); err != nil {
			err = fmt.Errorf("%w", err)
		}
	}
	return err
}

func (r runtimeRegistry) inc(model, metric string) error {
	key := model + "::" + metric
	if _, exists := r[key]; exists {
		r[key].inc()
		return nil
	}
	return fmt.Errorf("unable to find runtime metric %s", key)
}

func (r runtimeRegistry) incAll(models []string, metric string) error {
	var err error
	for _, modelName := range models {
		if err = r.inc(modelName, metric); err != nil {
			err = fmt.Errorf("%w", err)
		}
	}
	return err
}

func (r runtimeRegistry) now(model, metric string) error {
	return r.set(model, metric, float64(time.Now().UTC().Unix()))
}

func (r runtimeRegistry) nowAll(models []string, metric string) error {
	return r.setAll(models, metric, float64(time.Now().UTC().Unix()))
}
