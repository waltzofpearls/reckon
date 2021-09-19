package config

import (
	"errors"
	"strings"

	"github.com/waltzofpearls/reckon/logs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// WatchList can decode either a YAML config or a comma separated list
//
// - YAML config has a list of metric queries and model names for each everyone of them.
//   When using YAML, MODELS (config.Models) env var is not needed.
// - comma separated list doesn't specify model names for each metric query, models names
//   come from MODELS (config.Models) env var and it will be set with fillEmpty method.
//   When using comma separated list, MODELS env var is required
type WatchList struct {
	logger logs.Logger
	// map{
	//     "metric1": {"model1", "model2"},
	//     "metric2": {"model1"},
	//     "metric3": {"model2"},
	// }
	list map[string][]string
}

func newWatchList(lg logs.Logger) *WatchList {
	return &WatchList{
		logger: lg,
	}
}

// Decode either a YAML config or a comma separated list for backward compatibility
//
// WATCH_LIST comma separated list or inline yaml
// sensehat_temperature,sensehat_humidity,sensehat_pressure
// {sensehat_temperature: [Prophet, LSTM], sensehat_humidity: [Prophet], sensehat_pressure: [LSTM]}
func (w *WatchList) Decode(value string) error {
	if len(value) == 0 {
		return errors.New("WATCH_LIST cannot be empty")
	}
	w.list = make(map[string][]string)
	if err := yaml.Unmarshal([]byte(value), &w.list); err != nil {
		w.logger.Info("received comma separated list from WATCH_LIST", zap.String("value", value))
		watchList := strings.Split(value, ",")
		for _, metricQuery := range watchList {
			w.list[metricQuery] = []string{}
		}
	} else {
		w.logger.Info("received YAML config from WATCH_LIST", zap.String("value", value))
	}
	return nil
}

func (w *WatchList) List() map[string][]string {
	return w.list
}

func (w *WatchList) fillEmpty(models []string) {
	for metricQuery, modelNames := range w.list {
		if len(modelNames) == 0 {
			w.list[metricQuery] = models
		}
	}
}
