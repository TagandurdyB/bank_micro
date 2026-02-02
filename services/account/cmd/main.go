package main

import (
	"bank_micro/pkg/database"
	"bank_micro/services/account/config"
	"bank_micro/services/account/internal/app"
	"log"
)

func main() {
	config.LoadAccountConfig()

	db, err := database.ConnectPostgres(config.Cfg.DbURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	container, err := app.InitDependencies(db, config.Cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("InitDependencies error: %v", err)
	}
	grpcApp := app.NewGRPCServer(config.Cfg.GRPCAddr, container)

	log.Printf("Account Service starting on %s...", config.Cfg.GRPCAddr)
	if err := grpcApp.Run(); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
