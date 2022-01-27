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
#### Example Program - Loading program config from yaml file (main.go)
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

## Notes

Currently the only supported primitive types in config structs are:
- `string`
- `bool`
- `int`,`int8`,`int16` ,`int32`, `int64` 
- `float32`,`float64`

I plan on adding other types very shortly.

## Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b feature/new-feature`.
- Then _commit_ your changes `git commit -m 'Added new feature`
- Make a _push_ to your _branch_ `git push origin feature/new-feature`.
- Submit a **Pull Request** so that I can review your changes

