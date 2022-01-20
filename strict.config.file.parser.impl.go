package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	goenvloader "github.com/MartinSimango/go-envloader"
	"gopkg.in/yaml.v2"
)

type StrictConfigFileParser struct {
	StrictConfig      interface{}
	FileConfiguration *FileConfiguration
}

// Check we implement interface
var _ ConfigFileParser = &StrictConfigFileParser{}

func NewStrictConfigFileParser(strictConfig interface{},
	fileConfiguration *FileConfiguration) *StrictConfigFileParser {
	return &StrictConfigFileParser{
		StrictConfig:      strictConfig,
		FileConfiguration: fileConfiguration,
	}
}

func DefaultStrictYamlConfigFileParser(configFile string, config interface{}, strictConfig interface{}) *StrictConfigFileParser {
	fileConfig := &FileConfiguration{
		FileName:               configFile,
		FileFormat:             YAML,
		FileInputConfiguration: config,
		EnvironmentLoader:      goenvloader.NewBraceEnvironmentLoader(),
	}
	return NewStrictConfigFileParser(
		strictConfig,
		fileConfig,
	)
}

func DefaultStrictPropertyConfigFileParser(configFile string, config interface{}, strictConfig interface{}) *StrictConfigFileParser {
	fileConfig := &FileConfiguration{
		FileName:               configFile,
		FileFormat:             PROPERTY,
		FileInputConfiguration: config,
		EnvironmentLoader:      goenvloader.NewBraceEnvironmentLoader(),
	}

	return NewStrictConfigFileParser(
		strictConfig,
		fileConfig,
	)
}

func (cf *StrictConfigFileParser) ParseConfig() (interface{}, error) {
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

func (cf *StrictConfigFileParser) parseYamlFile() (interface{}, error) {
	yamlFile, err := ioutil.ReadFile(cf.FileConfiguration.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	if err = yaml.Unmarshal(yamlFile, cf.FileConfiguration.FileInputConfiguration); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	yamlFieldErrors := cf.extractEnv(cf.FileConfiguration.FileInputConfiguration, cf.StrictConfig)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  YAML,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.StrictConfig, nil

}

func (cf *StrictConfigFileParser) parsePropertiesFile() (interface{}, error) {
	if _, err := toml.DecodeFile(cf.FileConfiguration.FileName, cf.FileConfiguration.FileInputConfiguration); err != nil {
		return nil, err
	}

	yamlFieldErrors := cf.extractEnv(cf.FileConfiguration.FileInputConfiguration, cf.StrictConfig)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  PROPERTY,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.StrictConfig, nil

}

func (cf *StrictConfigFileParser) extractEnv(configuration interface{}, output interface{}) []FieldError {
	inputConfig := reflect.ValueOf(configuration).Elem()
	outputConfig := reflect.ValueOf(output).Elem()
	fieldErrors := []FieldError{}

	if inputConfig.NumField() != outputConfig.NumField() {
		panic(fmt.Sprintf("input and output configuration structs are not compatible\nInput Struct:%s\nOutput Struct:%s",
			inputConfig.Type().Name(),
			outputConfig.Type().Name(),
		))
	}

	for i := 0; i < inputConfig.NumField(); i++ {
		switch inputConfig.Field(i).Kind() {
		case reflect.String:
			err := cf.setOutputConfigValue(inputConfig, outputConfig, i)
			fieldErrors = AppendFieldError(fieldErrors, err)
		case reflect.Struct:
			err := cf.extractEnv(inputConfig.Field(i).Addr().Interface(), outputConfig.Field(i).Addr().Interface())
			fieldErrors = append(fieldErrors, err...)

		default:
			panic(fmt.Sprintf("types within input configuration struct needs to be string or struct:\nInput Struct: %s",
				inputConfig.Type().Name(),
			))

		}
	}
	return fieldErrors
}

func (cf *StrictConfigFileParser) setOutputConfigValue(inputConfig, outputConfig reflect.Value, index int) *FieldError {
	field := outputConfig.Field(index)
	switch field.Kind() {
	case reflect.String:
		value, err := cf.FileConfiguration.EnvironmentLoader.LoadStringFromEnv(inputConfig.Field(index).String())
		if err != nil {
			return &FieldError{
				Field:        inputConfig.Type().Field(index).Name + ": " + inputConfig.Field(index).String(),
				ErrorMessage: cf.getErrorMessage(err, inputConfig.Field(index)),
			}
		}
		outputConfig.Field(index).SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := cf.FileConfiguration.EnvironmentLoader.LoadIntFromEnv(inputConfig.Field(index).String())
		if err != nil {
			return &FieldError{
				Field:        inputConfig.Type().Field(index).Name + ": " + inputConfig.Field(index).String(),
				ErrorMessage: cf.getErrorMessage(err, inputConfig.Field(index)),
			}
		}
		outputConfig.Field(index).SetInt(int64(value))
	default:
		panic(fmt.Sprintf("Unsupported type found in application config struct.\nStruct Name: %s\nProperty Type: %s\n",
			outputConfig.Type().Name(),
			outputConfig.Field(index).Type().String()))

	}
	return nil
}

func (cf *StrictConfigFileParser) getErrorMessage(err error, field reflect.Value) string {
	if strings.TrimSpace(field.String()) == "" {
		return "field is empty"
	}
	return err.Error()
}
