package service

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/mysll/toolkit"
)

type Config struct {
	ID        uint16
	Name      string
	NatsUrl   string
	Depend    []Dependency
	Expose    bool
	Addr      string
	Port      int
	FPS       int
	Debug     bool
	DebugPort int
	Args      map[string]string
}

func (c *Config) LoadFromJson(f string) error {
	data, err := toolkit.ReadFile(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func (c *Config) LoadFromToml(f string) error {
	_, err := toml.DecodeFile(f, c)
	return err
}
