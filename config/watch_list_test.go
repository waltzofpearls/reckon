package config

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waltzofpearls/reckon/mocks"
)

func TestWatchListDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		mock     func(logger *mocks.Logger)
		given    string
		wantList map[string][]string
		wantErr  error
	}{
		{
			name: "empty value received",
			mock: func(logger *mocks.Logger) {
				logger.EXPECT().Info(gomock.Any(), gomock.Any()).Times(0)
			},
			given:   "",
			wantErr: errors.New("WATCH_LIST cannot be empty"),
		},
		{
			name: "decode yaml config",
			mock: func(logger *mocks.Logger) {
				logger.EXPECT().Info("received YAML config from WATCH_LIST", gomock.Any())
			},
			given: "{sensehat_temperature: [Prophet, LSTM], sensehat_humidity: [Prophet], sensehat_pressure: [LSTM]}",
			wantList: map[string][]string{
				"sensehat_temperature": {"Prophet", "LSTM"},
				"sensehat_humidity":    {"Prophet"},
				"sensehat_pressure":    {"LSTM"},
			},
		},
		{
			name: "decode comma separated list",
			mock: func(logger *mocks.Logger) {
				logger.EXPECT().Info("received comma separated list from WATCH_LIST", gomock.Any())
			},
			given: "sensehat_temperature,sensehat_humidity,sensehat_pressure",
			wantList: map[string][]string{
				"sensehat_temperature": {},
				"sensehat_humidity":    {},
				"sensehat_pressure":    {},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := mocks.NewLogger(ctrl)
			test.mock(logger)

			watchList := newWatchList(logger)
			err := watchList.Decode(test.given)
			if test.wantErr != nil {
				assert.EqualError(t, err, test.wantErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.wantList, watchList.list)
			}
		})
	}
}
