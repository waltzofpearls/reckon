package model

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/model/api"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	logger := p.logger.With(zap.String("metric_name", data.Name), zap.Any("metric_labels", data.Labels),
		zap.Int("data_length", len(data.Values)), zap.String("want_duration", duration.String()))
	logger.Info("train model with data")

	creds, err := p.config.GRPCClientCreds()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load gRPC client credentials")
	}
	retryBackoff := grpc_retry.BackoffLinear(p.config.GRPCClientRetryBackoff)
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(retryBackoff),
		grpc_retry.WithCodes(codes.Unavailable, codes.NotFound, codes.Aborted),
	}
	conn, err := grpc.Dial(p.config.GRPCServerAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to gRPC server")
	}
	defer conn.Close()

	logger.Info("start prophet training")
	start := time.Now()

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
	maxRetries := grpc_retry.WithMax(p.config.GRPCClientMaxRetries)
	reply, err := client.Prophet(ctx, &request, maxRetries)
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

	logger.Info("prophet training completed",
		zap.Stringer("elapsed", time.Since(start)), zap.Int("forecasts", len(forecasts)))

	return forecasts, nil
}
