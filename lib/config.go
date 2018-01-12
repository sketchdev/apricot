package lib

import (
	"io/ioutil"
	"path"

	"github.com/BurntSushi/toml"
)

// Configuration is a holder for various configuration data
type Configuration struct {
	Engine         string
	Folders        []string
	ConnectionFile string
}

func (c Configuration) ConnectionString() (string, error) {
	data, err := ioutil.ReadFile(c.ConnectionFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// NewConfiguration builds a default configuration
func NewConfiguration(engine, connectionFile string) Configuration {
	return Configuration{engine, []string{path.Join("migrations", "current")}, connectionFile}
}

// NewConfigurationFromFile parses a Toml file to build a Configuration
func NewConfigurationFromFile(name string) (Configuration, error) {
	var config Configuration
	if _, err := toml.DecodeFile(name, &config); err != nil {
		return Configuration{}, err
	}
	return config, nil
}
