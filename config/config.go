package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
)

type Conf struct {
	Options struct {
		DocDir   string `toml:"doc_path"`
		Cert     string `toml:"cert_path"`
		CertKey  string `toml:"cert_key_path"`
		BindPort string `toml:"bind_port"`
		UseTls   bool   `toml:"use_tls"`
		Refresh  bool   `toml:"refresh_docs"`
	} `toml:"options"`
}

func ParseConfig(path string) (*Conf, error) {
	var conf *Conf
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, fmt.Errorf("Unable to read config %s: %s", path, err)
	}
	log.Println("Got config from", path)
	return conf, nil
}
