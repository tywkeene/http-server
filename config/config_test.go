package config

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	const configPath = "../config.toml"
	const badPath = "./konfig.camel"

	conf, err := ParseConfig(badPath)
	if conf != nil || err == nil {
		t.Fatal("Recieved config from bogus config file", badPath)
	}

	conf, err = ParseConfig(configPath)
	//Ensure we can get a correct configuration from an existing config file
	if conf == nil || err != nil {
		t.Fatalf("Failed to get config from %s: %s", configPath, err)
	}

}
