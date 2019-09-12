package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TLSRootCA             string `envconfig:"TLS_ROOT_CA"`
	TLSServerCert         string `envconfig:"TLS_SERVER_CERT"`
	TLSServerKey          string `envconfig:"TLS_SERVER_KEY"`
	TLSClientCert         string `envconfig:"TLS_CLIENT_CERT"`
	TLSClientKey          string `envconfig:"TLS_CLIENT_KEY"`
	GRPCServerAddress     string `envconfig:"GRPC_SERVER_ADDRESS"`
	ModelTrainingInterval int64  `envconfig:"MODEL_TRAINING_INTERVAL" default:"10"`

	PromClientURL                string `envconfig:"PROM_CLIENT_URL"`
	PromClientTLSCA              string `envconfig:"PROM_CLIENT_TLS_CA"`
	PromClientTLSCert            string `envconfig:"PROM_CLIENT_TLS_CERT"`
	PromClientTLSKey             string `envconfig:"PROM_CLIENT_TLS_KEY"`
	PromClientInsecureSkipVerify bool   `envconfig:"PROM_CLIENT_INSECURE_SKIP_VERIFY"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	return envconfig.Process("", c)
}

func (c *Config) PromClientTLS() (*tls.Config, error) {
	if c.PromClientTLSCA == "" &&
		c.PromClientTLSKey == "" &&
		c.PromClientTLSCert == "" &&
		!c.PromClientInsecureSkipVerify {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.PromClientInsecureSkipVerify,
		Renegotiation:      tls.RenegotiateNever,
	}

	if c.PromClientTLSCA != "" {
		pool, err := makeCertPool([]string{c.PromClientTLSCA})
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = pool
	}

	if c.PromClientTLSCert != "" && c.PromClientTLSKey != "" {
		err := loadCertificate(
			tlsConfig,
			c.PromClientTLSCert,
			c.PromClientTLSKey,
		)
		if err != nil {
			return nil, err
		}
	}

	return tlsConfig, nil
}

func makeCertPool(certFiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, certFile := range certFiles {
		pem, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, fmt.Errorf(
				"could not read certificate %q: %v", certFile, err)
		}
		ok := pool.AppendCertsFromPEM(pem)
		if !ok {
			return nil, fmt.Errorf(
				"could not parse any PEM certificates %q: %v", certFile, err)
		}
	}
	return pool, nil
}

func loadCertificate(config *tls.Config, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf(
			"could not load keypair %s:%s: %v", certFile, keyFile, err)
	}

	config.Certificates = []tls.Certificate{cert}
	config.BuildNameToCertificate()
	return nil
}
