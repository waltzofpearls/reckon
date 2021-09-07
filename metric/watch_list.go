package metric

import (
	"context"

	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type WatchList struct {
	config *config.Config
	logger *zap.Logger
	client *prom.Client
	store  *Store
}

func NewWatchList(cf *config.Config, lg *zap.Logger, cl *prom.Client, st *Store) WatchList {
	return WatchList{
		config: cf,
		logger: lg,
		client: cl,
		store:  st,
	}
}

func (w WatchList) Build(ctx context.Context) error {
	for _, query := range w.config.WatchList {
		metrics, err := w.client.GetCurrentMetricValue(ctx, query)
		if err != nil {
			return err
		}
		for _, metric := range metrics {
			del := newDelegate(w.logger, w.config, w.client, metric)
			w.store.Save(del.key, del)
		}
	}
	return nil
}
