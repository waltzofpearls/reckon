package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForecastIsZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		given Forecast
		want  bool
	}{
		{
			name:  "zero forecast",
			given: Forecast{},
			want:  true,
		},
		{
			name: "not a zero forecast",
			given: Forecast{
				Timestamp: 1626048772,
				Values:    map[Column]float64{Yhat: 25.0},
			},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.given.IsZero()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestForecastsNearest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		givenForecasts Forecasts
		givenTimestamp int64
		want           Forecast
	}{
		{
			name:           "empty forecasts",
			givenForecasts: Forecasts{},
			givenTimestamp: 1626048772,
			want:           Forecast{},
		},
		{
			name: "only 1 item in forecasts",
			givenForecasts: Forecasts{
				{
					Timestamp: 1626048772,
					Values:    map[Column]float64{Yhat: 25.0},
				},
			},
			givenTimestamp: 1626049999,
			want: Forecast{
				Timestamp: 1626048772,
				Values:    map[Column]float64{Yhat: 25.0},
			},
		},
		{
			name: "multiple forecasts",
			givenForecasts: Forecasts{
				{
					Timestamp: 1626048772,
					Values:    map[Column]float64{Yhat: 25.0},
				},
				{
					Timestamp: 1626048782,
					Values:    map[Column]float64{Yhat: 26.0},
				},
				{
					Timestamp: 1626048792,
					Values:    map[Column]float64{Yhat: 27.0},
				},
				{
					Timestamp: 1626048802,
					Values:    map[Column]float64{Yhat: 28.0},
				},
			},
			givenTimestamp: 1626048780,
			want: Forecast{
				Timestamp: 1626048782,
				Values:    map[Column]float64{Yhat: 26.0},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got := test.givenForecasts.Nearest(test.givenTimestamp)
			assert.Equal(t, test.want, got)
		})
	}
}
