package model

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Trainer interface {
	Train(context.Context, prom.Metric, time.Duration) (Forecasts, error)
}

const TypeProphet = "Prophet"

func New(cf *config.Config, lg *zap.Logger, name string) (Trainer, error) {
	switch name {
	case TypeProphet:
		return NewProphet(cf, lg), nil
	}
	return nil, fmt.Errorf("%s is not a valid model type", name)
}

type Column string

const (
	Yhat      Column = "yhat"
	YhatUpper Column = "yhat_upper"
	YhatLower Column = "yhat_lower"
)

type Forecast struct {
	Timestamp int64
	Values    map[Column]float64
}

func (f Forecast) IsZero() bool {
	return f.Timestamp == 0 && len(f.Values) == 0
}

type Forecasts []Forecast

func (f Forecasts) Nearest(timestamp int64) Forecast {
	if len(f) == 0 {
		return Forecast{}
	}
	current := f[0]
	for _, forecast := range f {
		if math.Abs(float64(timestamp-forecast.Timestamp)) < math.Abs(float64(timestamp-current.Timestamp)) {
			current = forecast
		}
	}
	return current
}

type Server struct {
	logger *zap.Logger
	signal chan os.Signal
}

func NewServer(lg *zap.Logger, sig chan os.Signal) Server {
	return Server{
		logger: lg,
		signal: sig,
	}
}

func (s Server) Start(ctx context.Context) func() error {
	return func() error {
		s.logger.Info("starting python gRPC server...")
		cancelCtx, cancelFn := context.WithCancel(context.Background())
		subProc := exec.CommandContext(cancelCtx, "python", "model/server/main.py")
		subProc.Stdout = os.Stdout
		subProc.Stderr = os.Stderr
		go func() {
			<-ctx.Done()
			s.logger.Info("stopping python gRPC server...")
			subProc.Process.Signal(<-s.signal)
			time.AfterFunc(5*time.Second, cancelFn)
		}()
		if err := subProc.Run(); err != nil {
			switch err.(type) {
			case *exec.ExitError:
				s.logger.Info("stopped python gRPC server")
			default:
				s.logger.Error("failed to run python gRPC server", zap.Error(err))
				return err
			}
		}
		return nil
	}
}
