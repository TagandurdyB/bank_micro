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

	var accAddr, transAddr string
	
	if pkg.CoreCfg.IsDocker {
		accAddr = "account-service:" + gRPCPort1
		transAddr = "transaction-service:" + gRPCPort2
	} else {
		accAddr = host + ":" + gRPCPort1
		transAddr = host + ":" + gRPCPort2
	}

	Cfg = &TransactionConfig{
		CoreConfig:         pkg.CoreCfg,
		AccountServiceAddr: accAddr,
		GRPCAddr:           transAddr,
	}
}
