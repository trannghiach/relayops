package main

import (
	"context"
	"log"

	"relayops/internal/broker"
	"relayops/internal/config"
	"relayops/internal/db"
	workerpkg "relayops/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pg := db.NewPostgres(ctx, cfg.DatabaseURL)
	nc := broker.NewNATS(cfg.NATSURL)

	consumer := workerpkg.NewConsumer(pg, nc)

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("failed to start worker: %v", err)
	}

	log.Println("worker is listening for events...")

	select {}
}
