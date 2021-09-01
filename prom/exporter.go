package prom

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/waltzofpearls/reckon/config"
	"go.uber.org/zap"
)

type Exporter struct {
	config *config.Config
	logger *zap.Logger
	server *http.Server
}

func NewExporter(cf *config.Config, lg *zap.Logger) Exporter {
	return Exporter{
		config: cf,
		logger: lg,
	}
}

func RegisterExporterEndpoints() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusSeeOther)
	})
	http.Handle("/metrics", promhttp.Handler())
}

func (e *Exporter) Start(ctx context.Context) func() error {
	return func() error {
		e.server = &http.Server{
			Addr: e.config.PromExporterAddr,
		}
		addrField := zap.String("addr", e.config.PromExporterAddr)
		e.logger.Info("starting prometheus exporter server...", addrField)
		if err := e.server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("prometheus exporter server stopped: %w", err)
		}
		e.logger.Info("prometheus exporter server stopped", addrField)
		return nil
	}
}

func (e *Exporter) Shutdown(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		e.logger.Info("shutting down prometheus exporter server",
			zap.String("addr", e.config.PromExporterAddr))

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := e.server.Shutdown(timeout); err != nil {
			err = fmt.Errorf("failed to gracefully stop prometheus exporter server: %w", err)
		}
		return nil
	}
}
