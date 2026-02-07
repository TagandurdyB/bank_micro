package main

import (
	"bank_micro/pkg/database"
	"bank_micro/services/account/config"
	"log"
)

func main() {
	config.LoadAccountConfig()

	pgx, err := database.ConnectPgx(config.Cfg.DbURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	container, err := InitDependencies(*pgx, config.Cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("InitDependencies error: %v", err)
	}
	// grpcApp := app.NewGinServer(config.Cfg.GRPCAddr, container)

	// log.Printf("Account Service starting on %s...", config.Cfg.GRPCAddr)
	// if err := grpcApp.Run(); err != nil {
	// 	log.Fatalf("Failed to serve gRPC: %v", err)
	// }

	log.Printf("Account Service starting on %s...", config.Cfg.GRPCAddr)
	if err := container.Gin.Run(config.Cfg.GRPCAddr); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
