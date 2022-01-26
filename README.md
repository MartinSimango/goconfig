# goconfig
Simple configuration file parser for property and yaml files. Allows for fields in config files to be environment variables with default values.


## Install

```
go get github.com/MartinSimango/goconfig
```

## Example
Below is an example config yaml file that goconfig can parse:
#### Config File (app.yaml)
``` yaml
# app.yaml

port: ${SERVICE_PORT,8000}
service-name: service
db: 
  host: ${DB_HOST, 127.0.0.1}
  port: ${DB_PORT,8890}

```
The format of a config value is: 
``` yaml
configValue: ${ENVIRONMENT_VARIABLE,default_value}  
# OR 
configValue: ${ENVIRONMENT_VARIABLE} 
# OR 
configValue: value
```

Below is an example of the code that parse the config file and stores the config within a struct.
#### Example Program (main.go)
``` go
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

```

