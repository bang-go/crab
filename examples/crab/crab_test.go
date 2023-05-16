package crab_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/viperx"
	"log"
)

func baseSetting() {
	err := crab.Use(crab.UseAppLog(), crab.UseViper(&viperx.Config{ConfigFormat: viperx.FileFormatYaml, ConfigPaths: []string{"./"}, ConfigNames: []string{"app"}}))
	if err != nil {
		log.Fatal(err)
	}
}
