package lib

import (
	"path"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Engine     string
	Migrations []string
}

func NewConfiguration(engine string) Configuration {
	config := Configuration{}
	config.Engine = engine
	config.Migrations = []string{path.Join("migrations", "current")}
	return config
}

func NewConfigurationFromFile(name string) (Configuration, error) {
	var config Configuration
	if _, err := toml.DecodeFile(name, &config); err != nil {
		return Configuration{}, err
	}
	return config, nil
}
