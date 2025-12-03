package env_test

import (
	"log"
	"testing"

	"github.com/bang-go/crab/core/base/env"
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
