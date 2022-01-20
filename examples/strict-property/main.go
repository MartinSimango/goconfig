package main

import (
	"fmt"

	goenvloader "github.com/MartinSimango/go-envloader"
	"github.com/MartinSimango/goconfig"
)

type PropertyServiceConfiguration struct {
	Port        string
	ServiceName string
	DbHost      string
	DbPort      string
}

type StrictPropertyServiceConfiguration struct {
	Port        int
	ServiceName string
	DbHost      string
	DbPort      int
}

func main() {
	fileConfig := goconfig.NewFileConfiguration("app.properties", goconfig.PROPERTY, &PropertyServiceConfiguration{}, goenvloader.NewBraceEnvironmentLoader())
	fileParser := goconfig.NewStrictConfigFileParser(&StrictPropertyServiceConfiguration{}, fileConfig)

	config, err := fileParser.ParseConfig()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", config)
	}
}
