package service

import (
	"encoding/json"
	"fmt"
	"turboengine/core/api"

	"github.com/BurntSushi/toml"
	"github.com/mysll/toolkit"
)

type Config struct {
	Dev       bool
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
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return c.Valid()
}

func (c *Config) LoadFromToml(f string) error {
	_, err := toml.DecodeFile(f, c)
	if err != nil {
		return err
	}

	return c.Valid()
}

func (c *Config) Valid() error {
	if c.ID > uint16(api.MAX_SID) {
		return fmt.Errorf("service id must in 0 ~ %d", api.MAX_SID)
	}
	return nil
}
