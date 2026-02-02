package main

import (
	"bank_micro/pkg/database"
	"bank_micro/services/transaction/config"
	"bank_micro/services/transaction/internal/app"
	"log"
)

func main() {
	config.LoadTransactionConfig()

	db, err := database.ConnectPostgres(config.Cfg.DbURL)
	if err != nil {
		log.Fatalf("DB Error: %v", err)
	}

	container, err := app.InitDependencies(db, config.Cfg.RabbitMQURL, config.Cfg.AccountServiceAddr)
	if err != nil {
		log.Fatalf("DI Error: %v", err)
	}

	// Start RabbitMQ Worker in the background.
	container.Worker.Start()

	grpcApp := app.NewGRPCServer(config.Cfg.GRPCAddr, container)

	log.Printf("Transaction Service starting on %s...", config.Cfg.GRPCAddr)
	if err := grpcApp.Run(); err != nil {
		log.Fatalf("gRPC Serve Error: %v", err)
	}
}
