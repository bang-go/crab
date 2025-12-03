package aliyun_trace

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/bang-go/crab/core/base/tracex"
	"go.opentelemetry.io/otel/attribute"
	otlpTraceGrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

const (
	slsProjectHeader            = "x-sls-otel-project"
	slsInstanceIDHeader         = "x-sls-otel-instance-id"
	slsAccessKeyIDHeader        = "x-sls-otel-ak-id"
	slsAccessKeySecretHeader    = "x-sls-otel-ak-secret"
	slsSecurityTokenHeader      = "x-sls-otel-token"
	TraceExporterEndpointStdout = "stdout"
)

type SlsConfig struct {
	Project         string
	InstanceID      string
	AccessKeyID     string
	AccessKeySecret string
}
type Config struct {
	TraceExporterEndpoint         string
	TraceExporterEndpointInsecure bool
	ServiceName                   string
	ServiceNamespace              string
	ServiceVersion                string
	SlsConfig
	Resource           *resource.Resource
	ResourceAttributes map[string]string
	client             tracex.Client
}

func New(c *Config) (conf *Config, err error) {
	err = c.mergeResource()
	if err != nil {
		return
	}
	err = c.IsValid()
	if err != nil {
		return
	}
	conf = c
	return
}

func (c *Config) Start() (err error) {
	var traceExporter sdktrace.SpanExporter
	switch c.TraceExporterEndpoint {
	case TraceExporterEndpointStdout, "":
		traceExporter, err = tracex.NewExporterByStdout(stdouttrace.WithPrettyPrint())
	default:
		headers := map[string]string{}
		if c.Project != "" && c.InstanceID != "" {
			headers = map[string]string{
				slsProjectHeader:         c.Project,
				slsInstanceIDHeader:      c.InstanceID,
				slsAccessKeyIDHeader:     c.AccessKeyID,
				slsAccessKeySecretHeader: c.AccessKeySecret,
			}
		}
		// 使用GRPC方式导出数据
		traceSecureOption := otlpTraceGrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		traceExporter, err = tracex.NewExporterByOltpGrpc(context.Background(), c.TraceExporterEndpoint, traceSecureOption, otlpTraceGrpc.WithHeaders(headers), otlpTraceGrpc.WithCompressor(gzip.Name))
	}
	if err != nil {
		return err
	}

	c.client = tracex.New(&tracex.Config{Sampler: tracex.DefaultSampler})
	return c.client.Start(tracex.WithExporterOption(traceExporter), tracex.WithResourceOption(c.Resource))

}
func (c *Config) Stop() {
	c.client.Shutdown()
}

// 默认使用本机hostname作为hostname
func (c *Config) getDefaultResource() *resource.Resource {
	hostname, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(c.ServiceName),
		semconv.HostNameKey.String(hostname),
		semconv.ServiceNamespaceKey.String(c.ServiceNamespace),
		semconv.ServiceVersionKey.String(c.ServiceVersion),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.ProcessCommandKey.String(os.Args[0]),
	)
}

func (c *Config) mergeResource() (err error) {
	if c.Resource, err = resource.Merge(c.getDefaultResource(), c.Resource); err != nil {
		return
	}
	var keyValues []attribute.KeyValue
	for key, value := range c.ResourceAttributes {
		keyValues = append(keyValues, attribute.KeyValue{
			Key:   attribute.Key(key),
			Value: attribute.StringValue(value),
		})
	}
	newResource := resource.NewWithAttributes(semconv.SchemaURL, keyValues...)
	if c.Resource, err = resource.Merge(c.Resource, newResource); err != nil {
		return
	}
	return
}

// IsValid check config and return error if config invalid
func (c *Config) IsValid() error {
	if c.ServiceName == "" {
		return errors.New("empty service name")
	}
	if c.ServiceVersion == "" {
		return errors.New("empty service version")
	}
	if strings.Contains(c.TraceExporterEndpoint, "log.aliyuncs.com") && c.TraceExporterEndpointInsecure {
		return errors.New("insecure grpc is not allowed when send data to sls directly")
	}
	if strings.Contains(c.TraceExporterEndpoint, "log.aliyuncs.com") {
		if c.Project == "" || c.InstanceID == "" || c.AccessKeyID == "" || c.AccessKeySecret == "" {
			return errors.New("empty project, instanceID, accessKeyID or accessKeySecret when send data to sls directly")
		}
		if strings.ContainsAny(c.Project, "${}") ||
			strings.ContainsAny(c.InstanceID, "${}") ||
			strings.ContainsAny(c.AccessKeyID, "${}") ||
			strings.ContainsAny(c.AccessKeySecret, "${}") {
			return errors.New("invalid project, instanceID, accessKeyID or accessKeySecret when send data to sls directly, you should replace these parameters with actual values")
		}
	}
	return nil
}
