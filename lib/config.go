package lib

import (
	"path"

	"github.com/BurntSushi/toml"
)

// Configuration is a holder for various configuration data
type Configuration struct {
	Engine     string
	Migrations []string
}

// NewConfiguration builds a default configuration
func NewConfiguration(engine string) Configuration {
	config := Configuration{}
	config.Engine = engine
	config.Migrations = []string{path.Join("migrations", "current")}
	return config
}

// NewConfigurationFromFile parses a Toml file to build a Configuration
func NewConfigurationFromFile(name string) (Configuration, error) {
	var config Configuration
	if _, err := toml.DecodeFile(name, &config); err != nil {
		return Configuration{}, err
	}
	return config, nil
}
