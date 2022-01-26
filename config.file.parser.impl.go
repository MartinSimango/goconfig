package goconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	goenvloader "github.com/MartinSimango/go-envloader"
	dynamicstruct "github.com/Ompluscator/dynamic-struct"
	"gopkg.in/yaml.v2"
)

type ConfigFileParserImpl struct {
	FileConfiguration *FileConfiguration
}

// Check we implement interface
var _ ConfigFileParser = &ConfigFileParserImpl{}

// DefaulltConfigFileParser is a constructor
func NewConfigFileParser(fileConfiguration *FileConfiguration) *ConfigFileParserImpl {
	return &ConfigFileParserImpl{
		FileConfiguration: fileConfiguration,
	}
}

// ParseConfig parses the file configured the FileConfiguration struct variable's FileName parameter.
func (cf *ConfigFileParserImpl) ParseConfig() (interface{}, error) {
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

// convertToStringConfig converts all the primative fields of config into string types.
func convertToStringConfig(config interface{}) (interface{}, error) {
	inputConfig := reflect.ValueOf(config).Elem()
	stringConfig := dynamicstruct.NewStruct()

	for i := 0; i < inputConfig.NumField(); i++ {
		field := inputConfig.Field(i)
		fieldName := inputConfig.Type().Field(i).Name
		fieldTag := inputConfig.Type().Field(i).Tag
		switch field.Kind() {

		case reflect.Struct:
			internalStruct, err := convertToStringConfig(field.Addr().Interface())
			if err != nil {
				return nil, fmt.Errorf("failed to convert internal struct '%s',%w", fieldName, err)
			}
			stringConfig.AddField(fieldName, reflect.ValueOf(internalStruct).Elem().Interface(), string(fieldTag))

		default:
			stringConfig.AddField(fieldName, "", string(fieldTag))
		}
	}
	return stringConfig.Build().New(), nil

}
func (cf *ConfigFileParserImpl) parseYamlFile() (interface{}, error) {

	yamlFile, err := ioutil.ReadFile(cf.FileConfiguration.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	stringConfig, err := convertToStringConfig(cf.FileConfiguration.FileInputConfiguration)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(yamlFile, stringConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	yamlFieldErrors := setConfigurationValues(stringConfig, cf.FileConfiguration.FileInputConfiguration, cf.FileConfiguration.EnvironmentLoader)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  YAML,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.FileConfiguration.FileInputConfiguration, nil

}

func (cf *ConfigFileParserImpl) parsePropertiesFile() (interface{}, error) {

	stringConfig, err := convertToStringConfig(cf.FileConfiguration.FileInputConfiguration)
	if err != nil {
		return nil, err
	}

	if _, err := toml.DecodeFile(cf.FileConfiguration.FileName, stringConfig); err != nil {
		return nil, err
	}

	yamlFieldErrors := setConfigurationValues(stringConfig, cf.FileConfiguration.FileInputConfiguration, cf.FileConfiguration.EnvironmentLoader)

	if len(yamlFieldErrors) > 0 {
		return nil, &FileConfigError{
			ConfigFile:  cf.FileConfiguration.FileName,
			FileFormat:  PROPERTY,
			FieldErrors: yamlFieldErrors,
		}
	}
	return cf.FileConfiguration.FileInputConfiguration, nil

}

func setConfigurationValues(stringConfiguration interface{}, configuration interface{}, environmentLoader goenvloader.EnvironmentLoader) []FieldError {
	inputConfig := reflect.ValueOf(stringConfiguration).Elem()
	outputConfig := reflect.ValueOf(configuration).Elem()
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
			err := setSingularConfigValue(inputConfig, outputConfig, i, environmentLoader)
			fieldErrors = AppendFieldError(fieldErrors, err)
		case reflect.Struct:
			err := setConfigurationValues(inputConfig.Field(i).Addr().Interface(), outputConfig.Field(i).Addr().Interface(), environmentLoader)
			fieldErrors = append(fieldErrors, err...)
		default:
			panic(fmt.Sprintf("types within input configuration struct needs to be string or struct:\nInput Struct: %s",
				inputConfig.Type().Name(),
			))
		}
	}
	return fieldErrors
}

func setSingularConfigValue(inputConfig, outputConfig reflect.Value, valueStructIndex int,
	environmentLoader goenvloader.EnvironmentLoader) *FieldError {

	var err error
	field := outputConfig.Field(valueStructIndex)
	switch field.Kind() {
	case reflect.String:
		err = setValueForString(inputConfig.Field(valueStructIndex).String(), field, environmentLoader)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = setValueForInt(inputConfig.Field(valueStructIndex).String(), field, environmentLoader)
	default:
		panic(fmt.Sprintf("Unsupported type found in application config struct.\nStruct Name: %s\nProperty Type: %s\n",
			outputConfig.Type().Name(),
			outputConfig.Field(valueStructIndex).Type().String()))
	}

	if err != nil {
		return &FieldError{
			Field:        inputConfig.Type().Field(valueStructIndex).Name + ": " + inputConfig.Field(valueStructIndex).String(),
			ErrorMessage: getFieldErrorMessage(err, inputConfig.Field(valueStructIndex)),
		}
	}
	return nil
}

func setValueForString(input string, outputField reflect.Value,
	environmentLoader goenvloader.EnvironmentLoader) error {

	value, err := environmentLoader.LoadStringFromEnv(input)

	if err != nil {
		return err
	}

	outputField.SetString(value)

	return nil
}

func setValueForInt(input string, outputField reflect.Value,
	environmentLoader goenvloader.EnvironmentLoader) error {

	value, err := environmentLoader.LoadIntFromEnv(input)

	if err != nil {
		return err
	}

	outputField.SetInt(int64(value))

	return nil
}

func getFieldErrorMessage(err error, field reflect.Value) string {
	if strings.TrimSpace(field.String()) == "" {
		return "field is empty"
	}
	return err.Error()
}
