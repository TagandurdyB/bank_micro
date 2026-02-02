package config

import (
	pkg "bank_micro/pkg/config"
)

type TransactionConfig struct {
	*pkg.CoreConfig
	AccountServiceAddr string
	GRPCAddr           string
}

var Cfg *TransactionConfig

func LoadTransactionConfig() {
	pkg.LoadConfig("transaction")

	host := pkg.GetEnv("HOST", "0.0.0.0")
	gRPCPort1 := pkg.GetEnv("GRPC_PORT_1", "51051")
	gRPCPort2 := pkg.GetEnv("GRPC_PORT_2", "51052")

	var accAddr string
	if pkg.CoreCfg.IsDocker {
		accAddr = "account-service:" + gRPCPort1
	} else {
		accAddr = host + ":" + gRPCPort1
	}

	Cfg = &TransactionConfig{
		CoreConfig:         pkg.CoreCfg,
		AccountServiceAddr: accAddr,
		GRPCAddr:           host + ":" + gRPCPort2,
	}
}
