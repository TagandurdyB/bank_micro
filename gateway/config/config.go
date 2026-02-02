package config

import (
	pkg "bank_micro/pkg/config"
)

type GatewayConfig struct {
	*pkg.CoreConfig
	AccountServiceAddr     string
	TransactionServiceAddr string
	ServerAddr             string
}

var Cfg *GatewayConfig

func LoadGatewayConfig() {
	pkg.LoadConfig("")

	isDocker := pkg.CoreCfg.IsDocker

	host := pkg.GetEnv("HOST", "0.0.0.0")
	gRPCPort1 := pkg.GetEnv("GRPC_PORT_1", "51051")
	gRPCPort2 := pkg.GetEnv("GRPC_PORT_2", "51052")

	var accAddr, transAddr string
	if isDocker {
		accAddr = "account-service:" + gRPCPort1
		transAddr = "transaction-service:" + gRPCPort2
	} else {
		accAddr = host + ":" + gRPCPort1
		transAddr = host + ":" + gRPCPort2
	}

	var port string
	if isDocker {
		port = pkg.GetEnv("API_PORT", "8080")
	} else {
		port = pkg.GetEnv("LOCAL_API_PORT", "9080")
	}

	Cfg = &GatewayConfig{
		CoreConfig:             pkg.CoreCfg,
		AccountServiceAddr:     accAddr,
		TransactionServiceAddr: transAddr,
		ServerAddr:             host + ":" + port,
	}
}
