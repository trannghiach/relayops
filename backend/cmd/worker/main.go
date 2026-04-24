package main

import (
	"context"
	"log"

	"relayops/internal/broker"
	"relayops/internal/config"
	"relayops/internal/db"
	"relayops/internal/dispatcher"
	workerpkg "relayops/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pg := db.NewPostgres(ctx, cfg.DatabaseURL)
	nc := broker.NewNATS(cfg.NATSURL)

	consumer := workerpkg.NewConsumer(pg, nc)
	emailDispatcher := dispatcher.NewEmailDispatcher(cfg.DemoRandomFailure, cfg.DemoFailureRate)
	runner := workerpkg.NewRunner(pg, emailDispatcher, cfg.DemoFastBackoff)

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("failed to start worker: %v", err)
	}

	go runner.Start(ctx)

	log.Println("worker is running (consumer and runner)...")

	select {}
}
