package grpcx

import (
	"github.com/bang-go/crab/core/base/tracex/instrument/grpctrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
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

var defaultClientKeepaliveParams = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             2 * time.Second,  // wait 2 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

type ClientConfig struct {
	Addr        string
	Secure      bool
	Trace       bool
	TraceFilter grpctrace.Filter
}

type ClientCallFunc func(*grpc.ClientConn) (any, error)

type ClientEntity struct {
	*ClientConfig
	conn               *grpc.ClientConn
	dialOptions        []grpc.DialOption
	streamInterceptors []grpc.StreamClientInterceptor
	unaryInterceptors  []grpc.UnaryClientInterceptor
}

// TODO: metric, retry, load balance

func NewClient(conf *ClientConfig) Client {
	return &ClientEntity{
		ClientConfig:       conf,
		dialOptions:        []grpc.DialOption{},
		streamInterceptors: []grpc.StreamClientInterceptor{},
		unaryInterceptors:  []grpc.UnaryClientInterceptor{},
	}
}

func (c *ClientEntity) Dial() (conn *grpc.ClientConn, err error) {
	baseClientOption := []grpc.DialOption{grpc.WithKeepaliveParams(defaultClientKeepaliveParams)}
	if !c.Secure {
		baseClientOption = append(baseClientOption, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	//if c.Trace {
	//	traceOption := []grpc.DialOption{grpc.WithChainUnaryInterceptor(grpctrace.UnaryClientInterceptor()), grpc.WithChainStreamInterceptor(grpctrace.StreamClientInterceptor())}
	//	if c.TraceFilter != nil {
	//		traceOption = []grpc.DialOption{grpc.WithChainUnaryInterceptor(grpctrace.UnaryClientInterceptor(grpctrace.WithFilter(c.TraceFilter))), grpc.WithChainStreamInterceptor(grpctrace.StreamClientInterceptor(grpctrace.WithFilter(c.TraceFilter)))}
	//	}
	//	baseClientOption = append(baseClientOption, traceOption...)
	//}
	c.dialOptions = append(baseClientOption, c.dialOptions...)
	options := append(c.dialOptions, grpc.WithChainUnaryInterceptor(c.unaryInterceptors...), grpc.WithChainStreamInterceptor(c.streamInterceptors...))
	//c.conn, err = grpc.Dial(c.ClientConfig.Addr, options...)
	c.conn, err = grpc.NewClient(c.ClientConfig.Addr, options...)
	return c.conn, err
}

func (c *ClientEntity) DialWithCall(call ClientCallFunc) (any, error) {
	conn, err := c.Dial()
	if err != nil {
		return nil, err
	}
	return call(conn)
}

func (c *ClientEntity) Conn() *grpc.ClientConn {
	return c.conn
}
func (c *ClientEntity) Close() {
	_ = c.conn.Close()
}

func (c *ClientEntity) AddDailOptions(dialOption ...grpc.DialOption) {
	c.dialOptions = append(c.dialOptions, dialOption...)
}

func (c *ClientEntity) AddUnaryInterceptor(interceptor ...grpc.UnaryClientInterceptor) {
	c.unaryInterceptors = append(c.unaryInterceptors, interceptor...)
}

func (c *ClientEntity) AddStreamInterceptor(interceptor ...grpc.StreamClientInterceptor) {
	c.streamInterceptors = append(c.streamInterceptors, interceptor...)
}
