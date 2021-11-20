package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	Schedule string `envconfig:"SCHEDULE" default:"@every 120m"`
	Timezone string `envconfig:"TIMEZONE" default:"America/Vancouver"`

	PromExporterAddr string `envconfig:"PROM_EXPORTER_ADDR" default:":8080"`

	PromClientURL                string `envconfig:"PROM_CLIENT_URL"`
	PromClientTLSCA              string `envconfig:"PROM_CLIENT_TLS_CA"`
	PromClientTLSCert            string `envconfig:"PROM_CLIENT_TLS_CERT"`
	PromClientTLSKey             string `envconfig:"PROM_CLIENT_TLS_KEY"`
	PromClientInsecureSkipVerify bool   `envconfig:"PROM_CLIENT_INSECURE_SKIP_VERIFY"`

	WatchList *WatchList `envconfig:"WATCH_LIST"`
	Models    []string   `envconfig:"MODELS" default:"Prophet"`

	DefaultChunkSize time.Duration `envconfig:"DEFAULT_CHUNK_SIZE" default:"120m"`
	RollingWindow    time.Duration `envconfig:"ROLLING_WINDOW" default:"72h"`

	GRPCServerAddress string `envconfig:"GRPC_SERVER_ADDRESS" default:"localhost:18443"`
	GRPCRootCA        string `envconfig:"GRPC_ROOT_CA"`
	GRPCServerCert    string `envconfig:"GRPC_SERVER_CERT"`
	GRPCServerKey     string `envconfig:"GRPC_SERVER_KEY"`
	GRPCClientCert    string `envconfig:"GRPC_CLIENT_CERT"`
	GRPCClientKey     string `envconfig:"GRPC_CLIENT_KEY"`

	logger    *zap.Logger
	location  *time.Location
	chunkSize time.Duration
}

func New(lg *zap.Logger) *Config {
	return &Config{
		logger: lg,
	}
}

func (c *Config) Load() error {
	c.WatchList = newWatchList(c.logger)
	if err := envconfig.Process("", c); err != nil {
		return err
	}
	c.WatchList.fillEmpty(c.Models)
	return nil
}

func (c *Config) Location() *time.Location {
	if c.location == nil {
		var err error
		if c.location, err = time.LoadLocation(c.Timezone); err != nil {
			c.logger.Error("cannot parse timezone", zap.String("tz", c.Timezone), zap.Error(err))
			c.location = time.Local
		}
	}
	return c.location
}

func (c *Config) ChunkSize() time.Duration {
	if c.chunkSize == 0 {
		now := time.Now().In(c.Location())
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
		schedule, err := parser.Parse(c.Schedule)
		if err != nil {
			c.chunkSize = c.DefaultChunkSize
		} else {
			c.chunkSize = RoundUpDuration(schedule.Next(now).Sub(now), time.Minute)
		}
	}
	return c.chunkSize
}

func RoundUpDuration(toRound, roundOn time.Duration) time.Duration {
	toRound = toRound.Round(roundOn)
	if toRound == 0 {
		toRound += roundOn
	}
	return toRound
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

func (c *Config) GRPCClientCreds() (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(c.GRPCRootCA)) {
		return nil, errors.New("failed to append root CA cert")
	}
	certificate, err := tls.X509KeyPair([]byte(c.GRPCClientCert), []byte(c.GRPCClientKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed load client TLS key and cert")
	}
	return credentials.NewTLS(&tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	}), nil
}
