package main

import (
	"fmt"

	goenvloader "github.com/MartinSimango/go-envloader"
	"github.com/MartinSimango/goconfig"
)

type YamlServiceConfiguration struct {
	Port        string `yaml:"port"`
	ServiceName string `yaml:"service-name"`
	DB          struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"db"`
}

type StrictYamlServiceConfiguration struct {
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service-name"`
	DB          struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"db"`
}

func main() {
	fileConfig := goconfig.NewFileConfiguration("app.yaml", goconfig.YAML, &YamlServiceConfiguration{}, goenvloader.NewBraceEnvironmentLoader())
	fileParser := goconfig.NewStrictConfigFileParser(&StrictYamlServiceConfiguration{}, fileConfig)

	config, err := fileParser.ParseConfig()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", config)
	}
}
