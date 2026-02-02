package config

import (
	pkg "bank_micro/pkg/config"
)

type AccountConfig struct {
	*pkg.CoreConfig
	GRPCAddr string
}

var Cfg *AccountConfig

func LoadAccountConfig() {
	pkg.LoadConfig("account")

	host := pkg.GetEnv("HOST", "0.0.0.0")
	gRPCPort1 := pkg.GetEnv("GRPC_PORT_1", "51051")

	Cfg = &AccountConfig{
		CoreConfig: pkg.CoreCfg,
		GRPCAddr:   host + ":" + gRPCPort1,
	}
}
