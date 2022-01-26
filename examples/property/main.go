package main

import (
	"fmt"

	"github.com/MartinSimango/goconfig"
)

type PropertyServiceConfiguration struct {
	Port        int
	ServiceName string
	DbHost      string
	DbPort      int
}

func main() {
	fileConfig := goconfig.PropertyFileConfiguration("app.properties", &PropertyServiceConfiguration{})
	fileParser := goconfig.NewConfigFileParser(fileConfig)

	config, err := fileParser.ParseConfig()

	if err != nil {
		fmt.Println(err)
	} else {
		yamlConfig := config.(*PropertyServiceConfiguration) // cast if needed

		fmt.Printf("%+v", yamlConfig)
	}
}
