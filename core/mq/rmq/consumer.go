package rmq

import (
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
)

type Config = rmqClient.Config

// NewSimpleConsumer todo:待完善
func NewSimpleConsumer(config *Config, opts ...rmqClient.SimpleConsumerOption) (rmqClient.SimpleConsumer, error) {
	//rmqClient.WithZapLogger(logx.GetLogger())
	rmqClient.WithConnOptions()

	return rmqClient.NewSimpleConsumer(config, opts...)
}
