package app

import (
	proto "bank_micro/proto/gen"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	Server *grpc.Server
	Addr   string
}

func NewGRPCServer(addr string, container *Container) *GRPCServer {
	s := grpc.NewServer()

	proto.RegisterAccountServiceServer(s, container.Handler)

	return &GRPCServer{
		Server: s,
		Addr:   addr,
	}
}

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	return s.Server.Serve(lis)
}
