package crab_test

import (
	"log"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/viperx"
)

func baseSetting() {
	err := crab.Use(crab.UseAppLogx(), crab.UseViper(&viperx.Config{ConfigFormat: viperx.FileFormatYaml, ConfigPaths: []string{"./"}, ConfigNames: []string{"app"}}))
	if err != nil {
		log.Fatal(err)
	}
}
