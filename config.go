package main

import (
	// standard library
	"io/ioutil"

	// external packages
	"launchpad.net/goyaml"
)

// Config represents the application configuration structure (stored in the config.yaml).
type Config struct {
	DeviceId     string `yaml:"deviceId,omitempty"`
	AccessToken  string `yaml:"accessToken,omitempty"`
	SessionName  string `yaml:"sessionName,omitempty"`
	ClientId     string `yaml:"clientId,omitempty"`
	ClientSecret string `yaml:"clientSecret,omitempty"`
	CallbackUrl  string `yaml:"callbackUrl,omitempty"`
	HttpPort     string `yaml:"httpPort,omitempty"`
}

// CONFIG is the global accessible application configuration.
var CONFIG *Config = new(Config)

// readConfig reads the config.yaml and initializes CONFIG.
func readConfig() (err error) {
	file, err := ioutil.ReadFile("../src/github.com/FraBle/sparkdo/config.yaml")
	if err != nil {
		return
	}
	if err = goyaml.Unmarshal(file, CONFIG); err != nil {
		return
	}
	return
}
