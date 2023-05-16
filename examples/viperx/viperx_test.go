package viperx_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/viperx"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestViper(t *testing.T) {
	if err := viperx.Build(&viperx.Config{ConfigFormat: viperx.FileFormatYaml, ConfigPaths: []string{"./config"}, ConfigNames: []string{"a", "b"}}); err != nil {
		log.Fatal(err)
	}
	log.Println(viper.GetString("name"))
	log.Println(viper.GetString("mode"))
}

func TestCrabWithViper(t *testing.T) {
	crab.Build()
	if err := crab.Use(crab.UseViper(&viperx.Config{ConfigFormat: viperx.FileFormatYaml, ConfigPaths: []string{"./config"}, ConfigNames: []string{"a", "b"}})); err != nil {
		log.Fatal(err)
	}
	_ = crab.Start()
	defer crab.Close()
	log.Println(viper.GetString("name"))
	log.Println(viper.GetString("mode"))
}
