package goconfig

// ConfigFileParser parses the config file
type ConfigFileParser interface {
	ParseConfig() (interface{}, error)
}
