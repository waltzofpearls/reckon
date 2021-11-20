package model

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/model/api"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Prophet struct {
	config *config.Config
	logger *zap.Logger
}

func NewProphet(cf *config.Config, lg *zap.Logger) Prophet {
	return Prophet{
		config: cf,
		logger: lg,
	}
}

func (p Prophet) Train(ctx context.Context, data prom.Metric, duration time.Duration) (Forecasts, error) {
	// wait 1 second for the model server to start
	time.Sleep(time.Second)

	logger := p.logger.With(zap.String("metric_name", data.Name), zap.Any("metric_labels", data.Labels))
	logger.Info("train model with data", zap.Int("length", len(data.Values)), zap.String("duration", duration.String()))

	creds, err := p.config.GRPCClientCreds()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load gRPC client credentials")
	}
	conn, err := grpc.Dial(p.config.GRPCServerAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gRPC server")
	}
	defer conn.Close()

	client := api.NewForecastClient(conn)
	values := make([]*api.SamplePair, len(data.Values))
	for i, value := range data.Values {
		values[i] = &api.SamplePair{
			Timestamp: value.Time.Unix(),
			Value:     value.Value,
		}
	}
	request := api.ProphetRequest{
		Values:   values,
		Duration: duration.Minutes(),
	}
	reply, err := client.Prophet(ctx, &request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to cal gRPC method Prophet")
	}

	var forecasts Forecasts
	for _, value := range reply.Values {
		forecasts = append(forecasts, Forecast{
			Timestamp: int64(value.Timestamp),
			Values: map[Column]float64{
				Yhat:      value.Yhat,
				YhatLower: value.YhatLower,
				YhatUpper: value.YhatUpper,
			},
		})
	}

	return forecasts, nil
}
