package rmq

import (
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/bang-go/crab/core/base/logx"
)

type Config = rmqClient.Config

// NewSimpleConsumer todo:待完善
func NewSimpleConsumer(config *Config, opts ...rmqClient.SimpleConsumerOption) (rmqClient.SimpleConsumer, error) {
	rmqClient.WithZapLogger(logx.GetLogger())

	rmqClient.WithConnOptions()

	return rmqClient.NewSimpleConsumer(config, opts...)
}
