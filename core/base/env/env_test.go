package env_test

import (
	"github.com/bang-go/crab/core/base/env"
	"log"
	"testing"
)

func TestAppEnv(t *testing.T) {
	err := env.Build(env.WithAppEnv(env.PROD))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(env.IsProd())
	log.Println(env.AppEnv())
	_ = env.Build(env.WithAppKey("crab_app_env"))
}
