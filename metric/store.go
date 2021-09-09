package metric

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/waltzofpearls/reckon/logs"
	"go.uber.org/zap"
)

type Store struct {
	logger logs.Logger
	sync.Map
}

func NewStore(lg logs.Logger) *Store {
	return &Store{
		logger: lg,
	}
}

func (s *Store) Save(key string, del *delegate) {
	s.Map.Store(key, del)
}

func (s *Store) Get(key string) (del *delegate, err error) {
	value, ok := s.Map.Load(key)
	if !ok {
		err = fmt.Errorf("cannot find metric delegate %s", key)
		s.logger.Error("cannot find metric delegate", zap.String("metric_desc", key))
		return
	}
	del, ok = value.(*delegate)
	if !ok {
		err = fmt.Errorf("stored metric delegate %s is not *metric.delegate typed", key)
		s.logger.Error("stored metric delegate is not *metric.delegate typed",
			zap.String("metric_desc", key), zap.String("type", reflect.TypeOf(value).String()))
		return
	}
	return
}

func (s *Store) ForEach(fn func(key string, del *delegate)) {
	s.Map.Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			s.logger.Error("key is not string typed",
				zap.Any("key", k), zap.String("type", reflect.TypeOf(k).String()))
			return false
		}
		del, ok := v.(*delegate)
		if !ok {
			s.logger.Error("stored metric delegate is not *metric.delegate typed",
				zap.Any("key", k), zap.String("type", reflect.TypeOf(v).String()))
			return false
		}
		fn(key, del)
		return true
	})
}
