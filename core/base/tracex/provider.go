package tracex

import (
	"context"
	"errors"

	"github.com/bang-go/opt"
	"github.com/bang-go/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

func Provider() trace.TracerProvider {
	return otel.GetTracerProvider()
}

type Client interface {
	Start(opts ...opt.Option[options]) error
	Shutdown()
}

const (
	DefaultSampler float64 = 1.0
)

type ClientEntity struct {
	config   *Config
	ctx      context.Context
	tp       *sdktrace.TracerProvider
	shutdown func()
}

type Config struct {
	Sampler        float64 //采样率 默认1.0
	AppName        string
	ExporterConfig *ExporterConfig //option 选填 默认使用stdout
}
type ExporterConfig struct {
	Kind     uint
	Endpoint string
}

type options struct {
	providerOptions []sdktrace.TracerProviderOption
	sampler         sdktrace.Sampler
	resource        *resource.Resource
	exporter        sdktrace.SpanExporter
}

func New(conf *Config) Client {
	return &ClientEntity{
		config: conf,
		ctx:    context.Background(),
	}
}

func (c *ClientEntity) defaultResource() *resource.Resource {
	return resource.NewSchemaless(semconv.ServiceNameKey.String(c.config.AppName))
}

func (c *ClientEntity) defaultSampler() sdktrace.Sampler {
	sampler := util.If(c.config.Sampler > 0, c.config.Sampler, DefaultSampler)
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampler))
}

func (c *ClientEntity) defaultExporter() (sdktrace.SpanExporter, error) {
	return NewExporterByStdout()
}

func (c *ClientEntity) defaultOptions() (opts *options, err error) {
	exporter, err := c.defaultExporter()
	if err != nil {
		return
	}
	opts = &options{
		sampler:  c.defaultSampler(),
		resource: c.defaultResource(),
		exporter: exporter,
	}
	return
}

func (c *ClientEntity) Start(opts ...opt.Option[options]) error {
	o, err := c.defaultOptions()
	if err != nil {
		return err
	}
	opt.Each(o, opts...)
	var exporter sdktrace.SpanExporter
	if c.config.ExporterConfig != nil {
		if exporter, err = c.makeSimpleExporter(); err != nil {
			return err
		}
	}
	if o.exporter != nil {
		exporter = o.exporter
	}
	baseOptions := []sdktrace.TracerProviderOption{sdktrace.WithSampler(o.sampler), sdktrace.WithResource(o.resource), sdktrace.WithBatcher(exporter)}
	providerOptions := append(baseOptions, o.providerOptions...)
	c.tp = sdktrace.NewTracerProvider(providerOptions...)
	// 设置全局 TracerProvider，供整个应用使用
	otel.SetTracerProvider(c.tp)
	// 设置全局文本映射传播器，用于跨服务传播追踪上下文
	otel.SetTextMapPropagator(defaultPropagator())
	c.shutdown = func() {
		_ = c.tp.Shutdown(context.Background())
		_ = o.exporter.Shutdown(context.Background())
	}
	return nil
}
func (c *ClientEntity) Shutdown() {
	if c.shutdown != nil {
		c.shutdown()
	}
}

func WithProviderOption(providers ...sdktrace.TracerProviderOption) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.providerOptions = append(o.providerOptions, providers...)
	})
}

func WithSamplerOption(sampler sdktrace.Sampler) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.sampler = sampler
	})
}

func WithExporterOption(exporter sdktrace.SpanExporter) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.exporter = exporter
	})
}

func WithResourceOption(resource *resource.Resource) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.resource = resource
	})
}
func (c *ClientEntity) makeSimpleExporter() (sdktrace.SpanExporter, error) {
	switch c.config.ExporterConfig.Kind {
	case ExporterKindStdout:
		return NewExporterByStdout(stdouttrace.WithPrettyPrint())
	case ExporterKindOltpGrpc:
		return NewExporterByOltpGrpc(c.ctx, c.config.ExporterConfig.Endpoint)
	case ExporterKindOltpHttp:
		return NewExporterByOltpHttp(c.ctx, c.config.ExporterConfig.Endpoint)
	default:
		return nil, errors.New("unknown exporter kind")
	}

}
