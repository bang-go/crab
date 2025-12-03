package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/bang-go/opt"
)

const (
	PROD             = "prod"    //生产环境
	GRAY             = "gray"    //灰度环境
	PRE              = "pre"     //预发布环境
	TEST             = "test"    //测试环境
	DEV              = "dev"     //开发环境
	DefaultAppEnvKey = "APP_ENV" //默认应用环境变量
	DefaultEnv       = DEV       //默认环境
)

var appEnv string //应用环境变量

func Build(opts ...opt.Option[Options]) error {
	o := &Options{
		appKey: DefaultAppEnvKey,
	}
	opt.Each(o, opts...)
	appEnv = o.appEnv
	if appEnv == "" {
		appEnv = os.Getenv(o.appKey)
	}
	if appEnv == "" {
		appEnv = DefaultEnv
	}
	switch appEnv {
	case PROD, GRAY, PRE, TEST, DEV:
	default:
		return errors.New(fmt.Sprintf("Unknown environment variable: %s", appEnv))
	}
	return nil
}

func AppEnv() string {
	return appEnv
}

func IsProd() bool {
	return AppEnv() == PROD
}

func IsDev() bool {
	return AppEnv() == DEV
}

func IsTest() bool {
	return AppEnv() == TEST
}

func IsPre() bool {
	return AppEnv() == PRE
}

func IsGray() bool {
	return AppEnv() == GRAY
}
