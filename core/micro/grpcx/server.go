package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

type Server interface {
	AddServerOptions(serverOption ...grpc.ServerOption)
	AddUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor)
	AddStreamInterceptor(interceptor ...grpc.StreamServerInterceptor)
	Start(ServerRegisterFunc) error
}

type ServerWrapper struct {
	*ServerConfig
	serverOptions      []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

type ServerRegisterFunc func(*grpc.Server)
type ServerConfig struct {
	Addr string
}

// todo: 心跳检测，trace，metric

func NewServer(conf *ServerConfig) Server {
	return &ServerWrapper{
		ServerConfig:       conf,
		serverOptions:      nil,
		streamInterceptors: nil,
		unaryInterceptors:  nil,
	}
}

func (s *ServerWrapper) AddServerOptions(serverOption ...grpc.ServerOption) {
	s.serverOptions = append(s.serverOptions, serverOption...)
}

func (s *ServerWrapper) AddUnaryInterceptor(interceptor ...grpc.UnaryServerInterceptor) {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptor...)
}

func (s *ServerWrapper) AddStreamInterceptor(interceptor ...grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptor...)
}

func (s *ServerWrapper) Start(register ServerRegisterFunc) error {
	var err error
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	options := append(s.serverOptions, grpc.ChainUnaryInterceptor(s.unaryInterceptors...), grpc.ChainStreamInterceptor(s.streamInterceptors...))
	server := grpc.NewServer(options...)
	register(server)
	return server.Serve(lis)

}
