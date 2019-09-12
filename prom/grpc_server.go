package prom

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/waltzofpearls/reckon/api"
	"github.com/waltzofpearls/reckon/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCServer struct {
	tlsCert []byte
	tlsKey  []byte
	rootCA  []byte
	address string
	client  *Client
}

func NewGRPCServer(c *config.Config, client *Client) *GRPCServer {
	return &GRPCServer{
		tlsCert: []byte(c.TLSServerCert),
		tlsKey:  []byte(c.TLSServerKey),
		rootCA:  []byte(c.TLSRootCA),
		address: c.GRPCServerAddress,
		client:  client,
	}
}

func (g *GRPCServer) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("could not listen to %s %v", g.address, err)
	}

	serverOption, err := g.tlsServerOption()
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer(serverOption)
	api.RegisterMetricsServer(grpcServer, g)

	log.Println("Starting GRPC server", g.address)
	go func() {
		<-ctx.Done()
		log.Println("Gracefully stopping GRPC server")
		grpcServer.GracefulStop()
	}()
	return grpcServer.Serve(listen)
}

func (g *GRPCServer) tlsServerOption() (grpc.ServerOption, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(g.rootCA) {
		return nil, errors.New("failed to append root CA cert")
	}
	certificate, err := tls.X509KeyPair(g.tlsCert, g.tlsKey)
	if err != nil {
		return nil, fmt.Errorf("failed load server TLS key and cert: %s", err)
	}
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	return grpc.Creds(credentials.NewTLS(tlsConfig)), nil
}

func (g *GRPCServer) Query(ctx context.Context, req *api.QueryMetricsRequest) (*api.QueryMetricsResponse, error) {
	start, err := ptypes.Timestamp(req.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := ptypes.Timestamp(req.EndTime)
	if err != nil {
		return nil, err
	}
	metric, err := g.client.GetMetricRangeData(
		ctx, req.MetricName, req.Labels, start, end, 1*time.Hour, 0)
	if err != nil {
		return nil, err
	}
	return &api.QueryMetricsResponse{
		Metrics: []*api.Metric{
			&metric,
		},
	}, nil
}

func (g *GRPCServer) Write(ctx context.Context, req *api.WriteMetricsRequest) (*empty.Empty, error) {
	return nil, nil
}
