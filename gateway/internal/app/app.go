package app

import (
	"context"
	"bank_micro/proto/gen"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Container struct {
	Mux *runtime.ServeMux
}

func InitDependencies(ctx context.Context, accountAddr, transAddr string) (*Container, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Account Service Handler 
	err := gen.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, accountAddr, opts)
	if err != nil {
		return nil, err
	}

	// Transaction Service Handler
	err = gen.RegisterTransactionServiceHandlerFromEndpoint(ctx, mux, transAddr, opts)
	if err != nil {
		return nil, err
	}

	return &Container{
		Mux: mux,
	}, nil
}