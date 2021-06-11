package main

import (
	"github.com/waltzofpearls/reckon/agent"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := agent.Run(logger); err != nil {
		logger.Fatal("error running reckon", zap.Error(err))
	}
}
