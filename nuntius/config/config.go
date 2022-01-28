package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Payload struct {
	Jsonrpc string
	Id      interface{}
	Method  string
	Params  interface{}
	Url     string
}
type Config struct {
	cfg *Payload
}

func (c *Config) ParsePayload() (*Payload, error) {

	if _, err := toml.DecodeFile("config/config.toml", &c.cfg); err != nil {
		return nil, fmt.Errorf("decoding payload config file has been failed: %w", err)
	}
	fmt.Println(c.cfg)
	return c.cfg, nil
}
