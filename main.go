package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/model"
	"github.com/waltzofpearls/reckon/prom"
)

func main() {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		ctx, cancel := context.WithCancel(context.Background())

		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP,
			syscall.SIGTERM, syscall.SIGINT)
		go func() {
			sig := <-signals
			if sig == syscall.SIGHUP {
				log.Printf("Reloading config")
				<-reload
				reload <- true
			}
			cancel()
		}()

		err := run(ctx)
		if err != nil && err != context.Canceled {
			log.Fatal(err)
		}
	}
}

func run(ctx context.Context) error {
	conf := config.New()
	if err := conf.Load(); err != nil {
		return err
	}

	promClient := prom.NewClient(conf)
	grpcServer := prom.NewGRPCServer(conf, promClient)
	mlModel := model.New(conf)

	errChan := make(chan error)
	go func() {
		errChan <- grpcServer.Run(ctx)
	}()
	go func() {
		mlModel.Train(ctx)
		mlModel.UpdateOnInterval(ctx)
	}()

	return <-errChan
}
