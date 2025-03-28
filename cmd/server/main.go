package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/LuchnikKek/metrigo/internal/server"
	"github.com/LuchnikKek/metrigo/internal/storage"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context) error {
	cfg := NewConfig()
	cfg.ParseFlags()
	cfg.ParseEnvs()
	log.Printf("Config parsed: %+v\r\n", cfg)

	store := storage.NewInMemoryStorage()

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: server.MetricsRouter(store),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and serve: %v", err)
		}
	}()
	log.Printf("Listening on %s\r\n", srv.Addr)
	<-ctx.Done()

	log.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}
