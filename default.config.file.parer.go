package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type DefaultConfigFileParser struct {
	FileConfiguration *FileConfiguration
}

// Check we implement interface
var _ ConfigFileParser = &DefaultConfigFileParser{}

// NewDefaulltConfigFileParser is a constructor
func NewDefaulltConfigFileParser(fileConfiguration *FileConfiguration) *DefaultConfigFileParser {
	return &DefaultConfigFileParser{
		FileConfiguration: fileConfiguration,
	}
}

// ParseConfig parses the file configured the FileConfiguration struct variable's FileName parameter.
func (cf *DefaultConfigFileParser) ParseConfig() (interface{}, error) {
	_, err := os.Stat(cf.FileConfiguration.FileName)
	if err != nil {
		return nil, err
	}
	switch cf.FileConfiguration.FileFormat {
	case YAML:
		return cf.parseYamlFile()
	case JSON:
		return nil, fmt.Errorf("JSON parser not implemented yet")
	case PROPERTY:
		return cf.parsePropertiesFile()
	default:
		return nil, fmt.Errorf("unrecognized file format")
	}
}

func (cf *DefaultConfigFileParser) parseYamlFile() (interface{}, error) {
	yamlFile, err := ioutil.ReadFile(cf.FileConfiguration.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	if err = yaml.Unmarshal(yamlFile, cf.FileConfiguration.FileInputConfiguration); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	yamlFieldErrors := cf.extractEnv(cf.FileConfiguration.FileInputConfiguration)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  YAML,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.FileConfiguration.FileInputConfiguration, nil

}

func (cf *DefaultConfigFileParser) parsePropertiesFile() (interface{}, error) {
	if _, err := toml.DecodeFile(cf.FileConfiguration.FileName, cf.FileConfiguration.FileInputConfiguration); err != nil {
		return nil, err
	}

	yamlFieldErrors := cf.extractEnv(cf.FileConfiguration.FileInputConfiguration)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  PROPERTY,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.FileConfiguration.FileInputConfiguration, nil

}

func (cf *DefaultConfigFileParser) extractEnv(configuration interface{}) []FieldError {
	inputConfig := reflect.ValueOf(configuration).Elem()
	fieldErrors := []FieldError{}

	for i := 0; i < inputConfig.NumField(); i++ {
		field := inputConfig.Field(i)
		switch field.Kind() {
		case reflect.String:

			value, err := cf.FileConfiguration.EnvironmentLoader.LoadStringFromEnv(field.String())
			if err != nil {
				fieldErrors = AppendFieldError(fieldErrors, &FieldError{
					Field:        inputConfig.Type().Field(i).Name + ": " + field.String(),
					ErrorMessage: cf.getErrorMessage(err, field),
				})
			}
			field.SetString(value)
		case reflect.Struct:
			err := cf.extractEnv(field.Addr().Interface())
			fieldErrors = append(fieldErrors, err...)

		default:
			panic(fmt.Sprintf("types within input configuration struct needs to be string or struct:\nInput Struct: %s",
				inputConfig.Type().Name(),
			))

		}
	}
	return fieldErrors
}

func (cf *DefaultConfigFileParser) getErrorMessage(err error, field reflect.Value) string {
	if strings.TrimSpace(field.String()) == "" {
		return "field is empty"
	}
	return err.Error()
}
