package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Config vo.NacosClientParam

// NewClient 文档地址:https://github.com/nacos-group/nacos-sdk-go
func NewClient(conf *Config) (config_client.IConfigClient, error) {
	return clients.NewConfigClient(vo.NacosClientParam{ClientConfig: conf.ClientConfig, ServerConfigs: conf.ServerConfigs})
}
