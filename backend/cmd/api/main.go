package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"relayops/internal/broker"
	"relayops/internal/config"
	"relayops/internal/db"
	"relayops/internal/http/router"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pg := db.NewPostgres(ctx, cfg.DatabaseURL)
	nc := broker.NewNATS(cfg.NATSURL)

	r := router.NewRouter(pg, nc, cfg)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctxShut, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctxShut)
}
