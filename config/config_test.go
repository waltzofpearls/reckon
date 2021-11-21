package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConfigLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		envToSet   map[string]string
		wantError  error
		wantConfig *Config
	}{
		{
			name: "has error",
			envToSet: map[string]string{
				"WATCH_LIST":     "sensehat_pressure",
				"ROLLING_WINDOW": "not a duration",
			},
			wantError:  errors.New("envconfig.Process: assigning ROLLING_WINDOW to RollingWindow: converting 'not a duration' to type time.Duration. details: time: invalid duration \"not a duration\""),
			wantConfig: nil,
		},
		{
			name: "no error",
			envToSet: map[string]string{
				"WATCH_LIST":     "sensehat_pressure",
				"ROLLING_WINDOW": "5m",
			},
			wantError: nil,
			wantConfig: &Config{
				// include noop logger
				logger: zap.NewNop(),
				// from default values
				Schedule:          "@every 120m",
				Timezone:          "America/Vancouver",
				PromExporterAddr:  ":8080",
				Models:            []string{"Prophet"},
				DefaultChunkSize:  120 * time.Minute,
				GRPCServerAddress: "localhost:18443",
				// an env var set for the test case
				RollingWindow: 5 * time.Minute,
				WatchList: &WatchList{
					logger: zap.NewNop(),
					list: map[string][]string{
						"sensehat_pressure": {"Prophet"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envToSet != nil {
				for key, value := range test.envToSet {
					os.Setenv(key, value)
				}
				defer func() {
					for key := range test.envToSet {
						os.Unsetenv(key)
					}
				}()
			}
			c := New(zap.NewNop())
			err := c.Load()
			if test.wantError != nil {
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				assert.Equal(t, test.wantConfig, c)
			}
		})
	}
}

func TestConfigLocation(t *testing.T) {
	t.Parallel()

	newYork, _ := time.LoadLocation("America/New_York")

	tests := []struct {
		name          string
		givenTimezone string
		wantLocation  *time.Location
	}{
		{
			name:          "failed parsing timezone and fallback on local",
			givenTimezone: "NotAValidTimezone",
			wantLocation:  time.Local,
		},
		{
			name:          "can parse timezone and set location",
			givenTimezone: "America/New_York",
			wantLocation:  newYork,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := &Config{
				logger:   zap.NewNop(),
				Timezone: test.givenTimezone,
			}
			gotLocation := c.Location()
			assert.Equal(t, test.wantLocation, gotLocation)
			assert.Equal(t, test.wantLocation, c.location)
		})
	}
}

func TestConfigChunkSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		givenSchedule string
		wantDuration  time.Duration
	}{
		{
			name:          "failed parsing schedule and fallback on default chunk size",
			givenSchedule: "not a valid schedule",
			wantDuration:  120 * time.Minute,
		},
		{
			name:          "can parse schedule and set chuck size",
			givenSchedule: "@every 5m",
			wantDuration:  5 * time.Minute,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			c := &Config{
				Schedule:         test.givenSchedule,
				DefaultChunkSize: 120 * time.Minute,
			}
			gotDuration := c.ChunkSize()
			assert.Equal(t, test.wantDuration, gotDuration)
			assert.Equal(t, test.wantDuration, c.chunkSize)
		})
	}
}

func TestRoundUpDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		givenToRound time.Duration
		givenRoundOn time.Duration
		want         time.Duration
	}{
		{
			name:         "round up from seconds to one minute",
			givenToRound: 5 * time.Second,
			givenRoundOn: time.Minute,
			want:         time.Minute,
		},
		{
			name:         "round up minutes from minutes",
			givenToRound: 10 * time.Minute,
			givenRoundOn: time.Minute,
			want:         10 * time.Minute,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := RoundUpDuration(test.givenToRound, test.givenRoundOn)
			assert.Equal(t, test.want, got)
		})
	}
}
