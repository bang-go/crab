package viperx

import "github.com/spf13/viper"

const (
	FileFormatYaml = "yaml"
	FileFormatJson = "json"
)

type Config struct {
	ConfigFormat string   //配置文件的格式
	ConfigPaths  []string // 查找配置文件的路径
	ConfigNames  []string //配置文件名称(无扩展名)
}

func Build(conf *Config) error {
	var err error
	viper.SetConfigType(conf.ConfigFormat) // 查找配置文件的格式
	for _, value := range conf.ConfigPaths {
		viper.AddConfigPath(value) // 查找配置文件的路径
	}
	for _, value := range conf.ConfigNames {
		viper.SetConfigName(value)                   //配置文件名称(无扩展名)
		if err = viper.MergeInConfig(); err != nil { // 处理错误
			return err
		}
	}
	return nil
}
