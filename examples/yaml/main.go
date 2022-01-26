package main

import (
	"fmt"

	"github.com/MartinSimango/goconfig"
)

type YamlServiceConfiguration struct {
	Port        int    `yaml:"port"`
	ServiceName string `yaml:"service-name"`
	DB          struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"db"`
}

func main() {

	fileConfig := goconfig.YamlFileConfiguration("app.yaml", &YamlServiceConfiguration{})
	fileParser := goconfig.NewConfigFileParser(fileConfig)

	config, err := fileParser.ParseConfig()

	if err != nil {
		fmt.Println(err)
	} else {
		yamlConfig := config.(*YamlServiceConfiguration) // cast if needed

		fmt.Printf("%+v", yamlConfig)
	}
}
