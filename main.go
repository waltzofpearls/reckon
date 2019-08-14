package main

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/waltzofpearls/reckon/api"
	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	var server Server
	api.RegisterMetricsServer(grpcServer, server)

	address := ":3003"
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("could not listen to %s %v", address, err)
	}
	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listen))
}

type Server struct{}

func (Server) Query(ctx context.Context, req *api.QueryMetricsRequest) (*api.QueryMetricsResponse, error) {
	return &api.QueryMetricsResponse{
		Metrics: []*api.Metric{
			&api.Metric{
				Metric: map[string]string{"__name__": "foobar"},
				Value:  []float64{1},
			},
			&api.Metric{
				Metric: map[string]string{"__name__": "bazbuz"},
				Value:  []float64{2},
			},
		},
	}, nil
}

func (Server) Write(ctx context.Context, req *api.WriteMetricsRequest) (*empty.Empty, error) {
	return nil, nil
}
