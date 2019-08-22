package model

import (
	"context"
	"log"
	"time"

	"github.com/waltzofpearls/reckon/config"
)

type Model struct {
	interval time.Duration
}

func New(c *config.Config) *Model {
	return &Model{
		interval: time.Duration(c.ModelTrainingInterval) * time.Second,
	}
}

func (m *Model) Train(ctx context.Context) {
	log.Println("Train model")
}

func (m *Model) UpdateOnInterval(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	go func() {
		<-ctx.Done()
		log.Println("Stopping model training scheduler")
		ticker.Stop()
	}()
	for range ticker.C {
		log.Println("Schedule model training")
		m.Train(ctx)
	}
}
