package goconfig

import goenvloader "github.com/MartinSimango/go-envloader"

type FileFormat int

const (
	YAML FileFormat = iota
	JSON
	PROPERTY
)

type FileConfiguration struct {
	FileName               string
	FileFormat             FileFormat
	FileInputConfiguration interface{}
	EnvironmentLoader      goenvloader.EnvironmentLoader
}

func FileFormatToString(fileFormat FileFormat) string {
	switch fileFormat {
	case YAML:
		return "YAML"
	case PROPERTY:
		return "property"
	case JSON:
		return "JSON"
	}
	return "Unknown file format"
}

func NewFileConfiguration(fileName string,
	fileFormat FileFormat,
	fileInputConfiguration interface{},
	environmentLoader goenvloader.EnvironmentLoader) *FileConfiguration {
	return &FileConfiguration{
		FileName:               fileName,
		FileFormat:             fileFormat,
		FileInputConfiguration: fileInputConfiguration,
		EnvironmentLoader:      environmentLoader,
	}
}
