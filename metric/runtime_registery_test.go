package metric

import (
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeMetricLabels(t *testing.T) {
	t.Parallel()

	hostname, _ := os.Hostname()

	tests := []struct {
		name  string
		given runtimeMetric
		want  []string
	}{
		{
			name: "happy path",
			given: runtimeMetric{
				original: "test_original_metric",
				model:    "LSTM",
			},
			want: []string{"test_original_metric", "LSTM", hostname},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.given.labels()
			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestRuntimeMetricInto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given uint64
		want  float64
	}{
		{
			name:  "0",
			given: 0,
			want:  0,
		},
		{
			name:  "1",
			given: 4607182418800017408,
			want:  1,
		},
		{
			name:  "1.024",
			given: 4607290505191074300,
			want:  1.024,
		},
		{
			name:  "10",
			given: 4621819117588971520,
			want:  10,
		},
		{
			name:  "10.24",
			given: 4621954225577792635,
			want:  10.24,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rtMetric := runtimeMetric{
				valBits: test.given,
			}
			assert.Equal(t, test.want, rtMetric.into())
		})
	}
}

func TestRuntimeMetricSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given float64
		want  uint64
	}{
		{
			name:  "0",
			given: 0,
			want:  0,
		},
		{
			name:  "1",
			given: 1,
			want:  4607182418800017408,
		},
		{
			name:  "1.024",
			given: 1.024,
			want:  4607290505191074300,
		},
		{
			name:  "10",
			given: 10,
			want:  4621819117588971520,
		},
		{
			name:  "10.24",
			given: 10.24,
			want:  4621954225577792635,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rtMetric := new(runtimeMetric)
			rtMetric.set(test.given)
			assert.Equal(t, test.want, atomic.LoadUint64(&rtMetric.valBits))
		})
	}
}

func TestRuntimeMetricAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given float64
		want  uint64
	}{
		{
			name:  "1 + 0",
			given: 0,
			want:  4607182418800017408,
		},
		{
			name:  "1 + 1",
			given: 1,
			want:  4611686018427387904,
		},
		{
			name:  "1 + 1.024",
			given: 1.024,
			want:  4611740061622916350,
		},
		{
			name:  "1 + 10",
			given: 10,
			want:  4622382067542392832,
		},
		{
			name:  "1 + 10.24",
			given: 10.24,
			want:  4622517175531213947,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rtMetric := &runtimeMetric{
				valBits: 4607182418800017408, // == 1 in float64
			}
			rtMetric.add(test.given)
			assert.Equal(t, test.want, atomic.LoadUint64(&rtMetric.valBits))
		})
	}
}

func TestRuntimeMetricInc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		times int
	}{
		{
			name:  "inc 5 times",
			times: 5,
		},
		{
			name:  "inc 50 times",
			times: 50,
		},
		{
			name:  "inc 500 times",
			times: 500,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rtMetric := new(runtimeMetric)
			var wg sync.WaitGroup
			for i := 0; i < test.times; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					rtMetric.inc()
				}()
			}
			wg.Wait()
			assert.Equal(t, float64(test.times), rtMetric.into())
		})
	}
}

func TestNewRuntimeRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModels []string
		wantMapKeys []string
	}{
		{
			name:        "one model",
			givenModels: []string{"Prophet"},
			wantMapKeys: []string{
				"Prophet::reckon_prometheus_client_scrape_time_seconds",
				"Prophet::reckon_exporter_scraped_time_seconds",
				"Prophet::reckon_model_train_time_seconds",
				"Prophet::reckon_forecast_data_received_time_seconds",
				"Prophet::reckon_model_train_total",
				"Prophet::reckon_model_train_duration_seconds",
				"Prophet::reckon_model_train_errors_total",
				"Prophet::reckon_prometheus_client_scrape_errors_total",
				"Prophet::reckon_exporter_scraped_total",
				"Prophet::reckon_prometheus_client_scrape_total",
				"Prophet::reckon_data_scraped_duration_minutes",
				"Prophet::reckon_forecast_data_duration_minutes",
				"Prophet::reckon_training_data_duration_minutes",
				"Prophet::reckon_data_scraped_values",
				"Prophet::reckon_forecast_data_values",
				"Prophet::reckon_training_data_values",
			},
		},
		{
			name:        "two models",
			givenModels: []string{"Prophet", "LSTM"},
			wantMapKeys: []string{
				"Prophet::reckon_prometheus_client_scrape_time_seconds",
				"Prophet::reckon_exporter_scraped_time_seconds",
				"Prophet::reckon_model_train_time_seconds",
				"Prophet::reckon_forecast_data_received_time_seconds",
				"Prophet::reckon_model_train_total",
				"Prophet::reckon_model_train_duration_seconds",
				"Prophet::reckon_model_train_errors_total",
				"Prophet::reckon_prometheus_client_scrape_errors_total",
				"Prophet::reckon_exporter_scraped_total",
				"Prophet::reckon_prometheus_client_scrape_total",
				"Prophet::reckon_data_scraped_duration_minutes",
				"Prophet::reckon_forecast_data_duration_minutes",
				"Prophet::reckon_training_data_duration_minutes",
				"Prophet::reckon_data_scraped_values",
				"Prophet::reckon_forecast_data_values",
				"Prophet::reckon_training_data_values",
				"LSTM::reckon_prometheus_client_scrape_time_seconds",
				"LSTM::reckon_exporter_scraped_time_seconds",
				"LSTM::reckon_model_train_time_seconds",
				"LSTM::reckon_forecast_data_received_time_seconds",
				"LSTM::reckon_model_train_total",
				"LSTM::reckon_model_train_duration_seconds",
				"LSTM::reckon_model_train_errors_total",
				"LSTM::reckon_prometheus_client_scrape_errors_total",
				"LSTM::reckon_exporter_scraped_total",
				"LSTM::reckon_prometheus_client_scrape_total",
				"LSTM::reckon_data_scraped_duration_minutes",
				"LSTM::reckon_forecast_data_duration_minutes",
				"LSTM::reckon_training_data_duration_minutes",
				"LSTM::reckon_data_scraped_values",
				"LSTM::reckon_forecast_data_values",
				"LSTM::reckon_training_data_values",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotRegistry := newRuntimeRegistry("original_metric", test.givenModels)
			var gotMapKeys []string
			for key := range gotRegistry {
				gotMapKeys = append(gotMapKeys, key)
			}
			assert.ElementsMatch(t, test.wantMapKeys, gotMapKeys)
		})
	}
}

func TestRuntimeRegistrySet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModel  string
		givenMetric string
		givenValue  float64
		wantError   error
		wantValue   float64
	}{
		{
			name:        "runtime metric not exists",
			givenModel:  "notExists",
			givenMetric: "notExists",
			givenValue:  15,
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "happy path",
			givenModel:  "Prophet",
			givenMetric: "reckon_training_data_duration_minutes",
			givenValue:  15,
			wantValue:   15,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet"})
			err := reg.set(test.givenModel, test.givenMetric, test.givenValue)
			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)
				gotMetric, exists := reg[test.givenModel+"::"+test.givenMetric]
				require.True(t, exists)
				assert.Equal(t, test.wantValue, gotMetric.into())
			}
		})
	}
}

func TestRuntimeRegistrySetAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModels []string
		givenMetric string
		givenValue  float64
		wantError   error
		wantValue   float64
	}{
		{
			name:        "runtime metric not exists",
			givenModels: []string{"notExists"},
			givenMetric: "notExists",
			givenValue:  15,
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "set for one model",
			givenModels: []string{"Prophet"},
			givenMetric: "reckon_training_data_duration_minutes",
			givenValue:  15,
			wantValue:   15,
		},
		{
			name:        "set for two models",
			givenModels: []string{"Prophet", "LSTM"},
			givenMetric: "reckon_training_data_duration_minutes",
			givenValue:  25,
			wantValue:   25,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet", "LSTM"})
			err := reg.setAll(test.givenModels, test.givenMetric, test.givenValue)

			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)

				keysToSet := make(map[string]struct{})
				for _, modelName := range test.givenModels {
					keysToSet[modelName+"::"+test.givenMetric] = struct{}{}
				}
				for key, metric := range reg {
					if _, exists := keysToSet[key]; exists {
						assert.Equal(t, test.wantValue, metric.into())
					} else {
						assert.Equal(t, float64(0), metric.into())
					}
				}
			}
		})
	}
}

func TestRuntimeRegistryInc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModel  string
		givenMetric string
		wantError   error
	}{
		{
			name:        "runtime metric not exists",
			givenModel:  "notExists",
			givenMetric: "notExists",
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "happy path",
			givenModel:  "Prophet",
			givenMetric: "reckon_prometheus_client_scrape_errors_total",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet"})
			err := reg.inc(test.givenModel, test.givenMetric)
			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)
				gotMetric, exists := reg[test.givenModel+"::"+test.givenMetric]
				require.True(t, exists)
				assert.Equal(t, float64(1), gotMetric.into())
			}
		})
	}
}

func TestRuntimeRegistryIncAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModels []string
		givenMetric string
		wantError   error
	}{
		{
			name:        "runtime metric not exists",
			givenModels: []string{"notExists"},
			givenMetric: "notExists",
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "set for one model",
			givenModels: []string{"Prophet"},
			givenMetric: "reckon_prometheus_client_scrape_errors_total",
		},
		{
			name:        "set for two models",
			givenModels: []string{"Prophet", "LSTM"},
			givenMetric: "reckon_prometheus_client_scrape_errors_total",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet", "LSTM"})
			err := reg.incAll(test.givenModels, test.givenMetric)

			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)

				keysToSet := make(map[string]struct{})
				for _, modelName := range test.givenModels {
					keysToSet[modelName+"::"+test.givenMetric] = struct{}{}
				}
				for key, metric := range reg {
					if _, exists := keysToSet[key]; exists {
						assert.Equal(t, float64(1), metric.into())
					} else {
						assert.Equal(t, float64(0), metric.into())
					}
				}
			}
		})
	}
}

func TestRuntimeRegistryNow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModel  string
		givenMetric string
		wantError   error
	}{
		{
			name:        "runtime metric not exists",
			givenModel:  "notExists",
			givenMetric: "notExists",
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "happy path",
			givenModel:  "Prophet",
			givenMetric: "reckon_exporter_scraped_time_seconds",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet"})
			err := reg.now(test.givenModel, test.givenMetric)
			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)

				gotMetric, exists := reg[test.givenModel+"::"+test.givenMetric]
				require.True(t, exists)

				from := time.Now().UTC().Add(-2 * time.Second).Unix()
				to := time.Now().UTC().Add(2 * time.Second).Unix()
				assert.LessOrEqual(t, float64(from), gotMetric.into())
				assert.GreaterOrEqual(t, float64(to), gotMetric.into())
			}
		})
	}
}

func TestRuntimeRegistryNowAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		givenModels []string
		givenMetric string
		wantError   error
	}{
		{
			name:        "runtime metric not exists",
			givenModels: []string{"notExists"},
			givenMetric: "notExists",
			wantError:   errors.New("unable to find runtime metric notExists::notExists"),
		},
		{
			name:        "set for one model",
			givenModels: []string{"Prophet"},
			givenMetric: "reckon_exporter_scraped_time_seconds",
		},
		{
			name:        "set for two models",
			givenModels: []string{"Prophet", "LSTM"},
			givenMetric: "reckon_exporter_scraped_time_seconds",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			reg := newRuntimeRegistry("original", []string{"Prophet", "LSTM"})
			err := reg.nowAll(test.givenModels, test.givenMetric)

			if test.wantError != nil {
				require.Error(t, err)
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				require.NoError(t, err)

				keysToSet := make(map[string]struct{})
				for _, modelName := range test.givenModels {
					keysToSet[modelName+"::"+test.givenMetric] = struct{}{}
				}
				for key, metric := range reg {
					if _, exists := keysToSet[key]; exists {
						from := time.Now().UTC().Add(-2 * time.Second).Unix()
						to := time.Now().UTC().Add(2 * time.Second).Unix()
						assert.LessOrEqual(t, float64(from), metric.into())
						assert.GreaterOrEqual(t, float64(to), metric.into())
					} else {
						assert.Equal(t, float64(0), metric.into())
					}
				}
			}
		})
	}
}
