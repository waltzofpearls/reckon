package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/waltzofpearls/reckon/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	serverCert := os.Getenv("TLS_SERVER_CERT")
	serverKey := os.Getenv("TLS_SERVER_KEY")
	rootCA := os.Getenv("TLS_ROOT_CA")
	serverAddress := os.Getenv("GRPC_SERVER_ADDRESS")

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(rootCA)) {
		log.Fatal("failed to append client certs")
	}
	certificate, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	var server Server
	serverOption := grpc.Creds(credentials.NewTLS(tlsConfig))
	grpcServer := grpc.NewServer(serverOption)
	api.RegisterMetricsServer(grpcServer, server)

	log.Println("Server starting...")
	listen, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("could not listen to %s %v", serverAddress, err)
	}
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
