package main

import (
	"log"

	"github.com/waltzofpearls/reckon/agent"
	"github.com/waltzofpearls/reckon/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("an error occurred creating logger", err)
	}
	defer logger.Sync()

	build := config.NewBuildInfo(version, commit, date, goVersion, pythonVersion, goreleaserVersion)

	if err := agent.Run(logger, build); err != nil {
		logger.Fatal("an error occurred running reckon", zap.Error(err))
	}
}

// these build info variables will be set by ldflags during build time
var (
	version           = "0.0.0"
	commit            = "qwerty123456"
	date              = "0000-00-00T00:00:00Z"
	goVersion         = "0.0.0"
	pythonVersion     = "0.0.0"
	goreleaserVersion = "0.0.0"
)
