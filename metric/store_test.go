package metric

import (
	"errors"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/waltzofpearls/reckon/mocks"
	"go.uber.org/zap"
)

func TestStoreSave(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		toSave map[string]*Delegate
	}{
		{
			name: "save one delegate",
			toSave: map[string]*Delegate{
				"one": {key: "one"},
			},
		},
		{
			name: "concurrently save five delegate",
			toSave: map[string]*Delegate{
				"one":       {key: "one"},
				"another":   {key: "another"},
				"something": {key: "something"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store := NewStore(zap.NewNop())

			var wg sync.WaitGroup
			for key, delegate := range test.toSave {
				wg.Add(1)
				go func(key string, delegate *Delegate) {
					defer wg.Done()
					store.Save(key, delegate)

				}(key, delegate)
			}
			wg.Wait()

			saved := make(map[string]*Delegate)
			for key := range test.toSave {
				value, err := store.Get(key)
				require.NoError(t, err)
				saved[key] = value
			}

			assert.Equal(t, test.toSave, saved)
		})
	}
}

func TestStoreGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		store     *Store
		givenKey  string
		wantValue *Delegate
		wantError error
	}{
		{
			name: "key does not exist",
			store: func() *Store {
				store := NewStore(zap.NewNop())
				return store
			}(),
			givenKey:  "notExists",
			wantError: errors.New("cannot find metric delegate notExists"),
		},
		{
			name: "value is not a *Delegate",
			store: func() *Store {
				store := NewStore(zap.NewNop())
				store.Store("notADelegate", &Collector{})
				return store
			}(),
			givenKey:  "notADelegate",
			wantError: errors.New("stored metric delegate notADelegate is not *metric.Delegate typed"),
		},
		{
			name: "happy path",
			store: func() *Store {
				store := NewStore(zap.NewNop())
				store.Store("legitDelegate", &Delegate{key: "legitDelegate"})
				return store
			}(),
			givenKey:  "legitDelegate",
			wantValue: &Delegate{key: "legitDelegate"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			gotValue, err := test.store.Get(test.givenKey)
			if test.wantError != nil {
				assert.EqualError(t, err, test.wantError.Error())
			} else {
				assert.Equal(t, test.wantValue, gotValue)
			}
		})
	}
}

func TestStoreForEach(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mock func(logger *mocks.Logger) *Store
		want map[string]*Delegate
	}{
		{
			name: "key is not string typed",
			mock: func(logger *mocks.Logger) *Store {
				store := NewStore(logger)
				store.Map.Store(123, "error will be logged")
				logger.EXPECT().Error("key is not string typed", gomock.Any())
				return store
			},
			want: map[string]*Delegate{},
		},
		{
			name: "value is not *Delegate typed",
			mock: func(logger *mocks.Logger) *Store {
				store := NewStore(logger)
				store.Map.Store("anotherError", "error will be logged")
				logger.EXPECT().Error("stored metric delegate is not *metric.Delegate typed", gomock.Any())
				return store
			},
			want: map[string]*Delegate{},
		},
		{
			name: "happy path",
			mock: func(logger *mocks.Logger) *Store {
				store := NewStore(logger)
				store.Save("one", &Delegate{key: "one"})
				store.Save("another", &Delegate{key: "another"})
				return store
			},
			want: map[string]*Delegate{
				"one":     {key: "one"},
				"another": {key: "another"},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := mocks.NewLogger(ctrl)
			store := test.mock(logger)

			got := make(map[string]*Delegate)
			store.ForEach(func(key string, delegate *Delegate) {
				got[key] = delegate
			})
			assert.Equal(t, test.want, got)
		})
	}
}
