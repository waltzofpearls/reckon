package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	python3 "github.com/go-python/cpy3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/metric"
	"github.com/waltzofpearls/reckon/model"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func Run(lg *zap.Logger, info *config.BuildInfo) error {
	module, cleanup, err := model.InitPythonModule(lg)
	if err != nil {
		return err
	}
	defer cleanup()
	prom.RegisterExporterEndpoints()
	return reloadLoop(lg, info, module)
}

func reloadLoop(lg *zap.Logger, info *config.BuildInfo, module *python3.PyObject) error {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		ctx, cancel := context.WithCancel(context.Background())

		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP,
			syscall.SIGTERM, syscall.SIGINT)
		go func() {
			sig := <-signals
			if sig == syscall.SIGHUP {
				lg.Info("reloading config")
				<-reload
				reload <- true
			}
			cancel()
		}()

		err := runSchedulerAndExporter(ctx, lg, info, module)
		if err != nil && err != context.Canceled {
			return err
		}
	}
	return nil
}

func runSchedulerAndExporter(ctx context.Context, lg *zap.Logger, info *config.BuildInfo, module *python3.PyObject) error {
	config := config.New(lg)
	if err := config.Load(); err != nil {
		return err
	}
	client := prom.NewClient(config, lg)
	store := metric.NewStore(lg)
	watchlist := metric.NewWatchList(config, lg, client, store)
	scheduler := metric.NewScheduler(config, lg, client, store)
	collector := metric.NewCollector(config, lg, info, store)
	exporter := prom.NewExporter(config, lg)

	if err := watchlist.Build(ctx); err != nil {
		return err
	}
	if err := prometheus.Register(collector); err != nil {
		return err
	}
	defer prometheus.Unregister(collector)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(scheduler.Start(ctx, module))
	g.Go(exporter.Start(ctx))
	g.Go(scheduler.Shutdown(ctx))
	g.Go(exporter.Shutdown(ctx))

	return g.Wait()
}
