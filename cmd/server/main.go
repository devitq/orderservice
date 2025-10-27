package main

import (
	"log"

	"orderservice/internal/config"
	"orderservice/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.New(cfg)
	srv.RegisterServices()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server is running on port %d", cfg.GRPCPort)

	select {}
}
