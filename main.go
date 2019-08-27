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
	c := config.New()
	if err := c.Load(); err != nil {
		return err
	}

	g := prom.NewGRPCServer(c)
	m := model.New(c)

	errChan := make(chan error)
	go func() {
		errChan <- g.Run(ctx)
	}()
	go func() {
		m.Train(ctx)
		m.UpdateOnInterval(ctx)
	}()

	return <-errChan
}
