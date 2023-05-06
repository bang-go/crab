package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	AddDailOptions(...grpc.DialOption)
	AddUnaryInterceptor(interceptor ...grpc.UnaryClientInterceptor)
	AddStreamInterceptor(interceptor ...grpc.StreamClientInterceptor)
	Dial() (*grpc.ClientConn, error)
	DialWithCall(ClientCallFunc) (any, error)
	Conn() *grpc.ClientConn
	Close()
}

type ClientConfig struct {
	Addr   string
	Secure bool
}

type ClientCallFunc func(*grpc.ClientConn) (any, error)

type ClientWrapper struct {
	*ClientConfig
	conn               *grpc.ClientConn
	dialOptions        []grpc.DialOption
	streamInterceptors []grpc.StreamClientInterceptor
	unaryInterceptors  []grpc.UnaryClientInterceptor
}

// todo: 心跳检测，trace，metric

func NewClient(conf *ClientConfig) Client {
	return &ClientWrapper{
		ClientConfig:       conf,
		dialOptions:        []grpc.DialOption{},
		streamInterceptors: []grpc.StreamClientInterceptor{},
		unaryInterceptors:  []grpc.UnaryClientInterceptor{},
	}
}

func (c *ClientWrapper) Dial() (conn *grpc.ClientConn, err error) {
	options := append(c.dialOptions, grpc.WithChainUnaryInterceptor(c.unaryInterceptors...), grpc.WithChainStreamInterceptor(c.streamInterceptors...))
	if !c.Secure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	c.conn, err = grpc.Dial(c.ClientConfig.Addr, options...)
	return c.conn, err
}

func (c *ClientWrapper) DialWithCall(call ClientCallFunc) (any, error) {
	conn, err := c.Dial()
	if err != nil {
		return nil, err
	}
	return call(conn)
}

func (c *ClientWrapper) Conn() *grpc.ClientConn {
	return c.conn
}
func (c *ClientWrapper) Close() {
	_ = c.conn.Close()
}

func (c *ClientWrapper) AddDailOptions(dialOption ...grpc.DialOption) {
	c.dialOptions = append(c.dialOptions, dialOption...)
}

func (c *ClientWrapper) AddUnaryInterceptor(interceptor ...grpc.UnaryClientInterceptor) {
	c.unaryInterceptors = append(c.unaryInterceptors, interceptor...)
}

func (c *ClientWrapper) AddStreamInterceptor(interceptor ...grpc.StreamClientInterceptor) {
	c.streamInterceptors = append(c.streamInterceptors, interceptor...)
}
