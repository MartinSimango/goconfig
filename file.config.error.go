package goconfig

import "fmt"

type ConfigError interface {
	Error
}

type FileConfigError struct {
	ConfigFile  string
	FileFormat  FileFormat
	FieldErrors []FieldError
}

var _ ConfigError = &FileConfigError{}

func (fce *FileConfigError) Error() string {

	return fmt.Sprintf("found %d error(s) while parsing %s config file(%s):\n%s",
		len(fce.FieldErrors),
		FileFormatToString(fce.FileFormat),
		fce.ConfigFile,
		FieldErrorsToString(fce.FieldErrors))
}
