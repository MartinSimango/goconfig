package main

import (
	"fmt"

	"github.com/MartinSimango/goconfig"
)

type PropertyServiceConfiguration struct {
	Port        string
	ServiceName string
	DbHost      string
	DbPort      string
}

func main() {
	fileConfig := goconfig.PropertyFileConfiguration("app.properties", &PropertyServiceConfiguration{})
	fileParser := goconfig.NewDefaultConfigFileParser(fileConfig)

	config, err := fileParser.ParseConfig()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", config)
	}
}
