package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"hookwatch/internal/server"
)

func main() {
	cfg := server.DefaultConfig()

	flag.IntVar(&cfg.Port, "port", cfg.Port, "port to listen on")
	flag.IntVar(&cfg.MaxItems, "max-items", cfg.MaxItems, "max webhook requests to keep in memory")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(cfg)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}

	log.Println("hookwatch stopped")
}
