package metric

import (
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/zap"
)

type Store struct {
	logger *zap.Logger
	sync.Map
}

func NewStore(lg *zap.Logger) *Store {
	return &Store{
		logger: lg,
	}
}

func (s *Store) Save(key string, delegate *Delegate) {
	s.Map.Store(key, delegate)
}

func (s *Store) Get(key string) (delegate *Delegate, err error) {
	value, ok := s.Map.Load(key)
	if !ok {
		err = fmt.Errorf("cannot find metric delegate %s", key)
		s.logger.Error("cannot find metric delegate", zap.String("metric_desc", key))
		return
	}
	delegate, ok = value.(*Delegate)
	if !ok {
		err = fmt.Errorf("stored metric delegate %s is not *metric.Delegate typed", key)
		s.logger.Error("stored metric delegate is not *metric.Delegate typed",
			zap.String("metric_desc", key), zap.String("type", reflect.TypeOf(value).String()))
		return
	}
	return
}

func (s *Store) ForEach(fn func(key string, delegate *Delegate)) {
	s.Map.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			s.logger.Error("key is not string typed",
				zap.Any("key", k), zap.String("type", reflect.TypeOf(k).String()))
			return false
		}
		delegate, ok := v.(*Delegate)
		if !ok {
			s.logger.Error("stored metric delegate is not :metric.Delegate typed",
				zap.Any("key", k), zap.String("type", reflect.TypeOf(v).String()))
			return false
		}
		fn(key, delegate)
		return true
	})
}
