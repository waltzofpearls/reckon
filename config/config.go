package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	TLSRootCA             string `envconfig:"TLS_ROOT_CA"`
	TLSServerCert         string `envconfig:"TLS_SERVER_CERT"`
	TLSServerKey          string `envconfig:"TLS_SERVER_KEY"`
	TLSClientCert         string `envconfig:"TLS_CLIENT_CERT"`
	TLSClientKey          string `envconfig:"TLS_CLIENT_KEY"`
	GRPCServerAddress     string `envconfig:"GRPC_SERVER_ADDRESS"`
	ModelTrainingInterval int64  `envconfig:"MODEL_TRAINING_INTERVAL" default:"10"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	return envconfig.Process("", c)
}
