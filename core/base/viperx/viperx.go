package viperx

import (
	"bytes"
	"embed"
	"io/fs"

	"github.com/spf13/viper"
)

const (
	FileFormatYaml = "yaml"
	FileFormatJson = "json"
)

type Config struct {
	ConfigFormat string   // 配置文件的格式
	ConfigPaths  []string // 查找配置文件的路径
	ConfigNames  []string // 配置文件名称(无扩展名)

	// 新增：支持直接传入配置数据
	ConfigData []byte // 配置文件的原始数据（优先级最高）

	// 新增：支持 embed.FS
	ConfigFS   *embed.FS // 嵌入的文件系统
	ConfigFile string    // embed.FS 中的文件路径（例如："configs/env.yaml"）
}

func Build(conf *Config) error {
	var err error
	viper.SetConfigType(conf.ConfigFormat) // 设置配置文件格式

	// 优先级1: 使用直接传入的配置数据（最高优先级）
	if len(conf.ConfigData) > 0 {
		return viper.MergeConfig(bytes.NewReader(conf.ConfigData))
	}

	// 优先级2: 使用 embed.FS
	if conf.ConfigFS != nil && conf.ConfigFile != "" {
		data, err := fs.ReadFile(*conf.ConfigFS, conf.ConfigFile)
		if err != nil {
			return err
		}
		return viper.MergeConfig(bytes.NewReader(data))
	}

	// 优先级3: 使用传统的文件路径方式
	for _, value := range conf.ConfigPaths {
		viper.AddConfigPath(value) // 查找配置文件的路径
	}
	for _, value := range conf.ConfigNames {
		viper.SetConfigName(value)                   // 配置文件名称(无扩展名)
		if err = viper.MergeInConfig(); err != nil { // 处理错误
			return err
		}
	}
	return nil
}
