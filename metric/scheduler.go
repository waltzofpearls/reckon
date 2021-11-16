package metric

import (
	"context"
	"log"
	"sync"

	python3 "github.com/go-python/cpy3"
	"github.com/robfig/cron/v3"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Scheduler struct {
	config *config.Config
	logger *zap.Logger
	store  *Store

	sync.Mutex
	cron *cron.Cron
}

func NewScheduler(cf *config.Config, lg *zap.Logger, cl *prom.Client, st *Store) *Scheduler {
	return &Scheduler{
		config: cf,
		logger: lg,
		store:  st,
	}
}

func (s *Scheduler) Start(ctx context.Context, module *python3.PyObject) func() error {
	return func() error {
		s.store.ForEach(func(key string, del *delegate) {
			s.logger.Info("schedule initial model training",
				zap.String("metric", key), zap.Strings("models", del.modelNames))
			go func() {
				del.train(ctx, module)
				s.store.Save(key, del)
			}()
		})

		s.Mutex.Lock()
		log.Println("Scheduler::Start")
		s.cron = cron.New(cron.WithSeconds(), cron.WithLocation(s.config.Location()))
		s.Mutex.Unlock()

		if _, err := s.cron.AddFunc(s.config.Schedule, func() {
			s.store.ForEach(func(key string, del *delegate) {
				s.logger.Info("schedule subsequent model training",
					zap.String("metric", key), zap.Strings("models", del.modelNames))
				go func() {
					del.train(ctx, module)
					s.store.Save(key, del)
				}()
			})
		}); err != nil {
			return err
		}

		s.cron.Run()
		s.logger.Info("model training scheduler stopped")
		return nil
	}
}

func (s *Scheduler) Shutdown(ctx context.Context) func() error {
	return func() error {
		<-ctx.Done()
		s.logger.Info("stopping model training scheduler")

		s.Mutex.Lock()
		stop := s.cron.Stop()
		s.Mutex.Unlock()

		<-stop.Done()
		return stop.Err()
	}
}
