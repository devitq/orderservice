package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"orderservice/internal/config"
	"orderservice/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	srv := server.New(cfg)
	srv.RegisterServices()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	srv.Stop()
	log.Println("server stopped")
}
