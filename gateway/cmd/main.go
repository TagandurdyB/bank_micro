package main

import (
	"bank_micro/gateway/config"
	"bank_micro/gateway/internal/app"
	"context"
	"log"
	"net/http"
)

func main() {
	config.LoadGatewayConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Dependency Injection / Container Starting
	container, err := app.InitDependencies(ctx, config.Cfg.AccountServiceAddr, config.Cfg.TransactionServiceAddr)
	if err != nil {
		log.Fatalf("Failed to initialize Gateway: %v", err)
	}

	log.Printf("API Gateway (REST) starting on :%s", config.Cfg.ServerAddr)
	if err := http.ListenAndServe(config.Cfg.ServerAddr, container.Mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
