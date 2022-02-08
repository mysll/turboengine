package config

import (
	"turboengine/core/api"
	"turboengine/core/plugin"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	Name = "Configuration"
)

type Configuration struct {
	srv    api.Service
	client *consulapi.Client
}

func (c *Configuration) Prepare(srv api.Service, args ...any) {
	c.srv = srv
	c.client = args[0].(*consulapi.Client)
}

func (c *Configuration) Run() {

}

func (c *Configuration) Shut(api.Service) {
}

func (c *Configuration) Handle(cmd string, args ...any) any {
	return nil
}

func (c *Configuration) StoreKV(key string, value []byte) error {
	kv := &consulapi.KVPair{
		Key:   key,
		Flags: 0,
		Value: value,
	}
	_, err := c.client.KV().Put(kv, nil)
	return err
}

func (c *Configuration) GetKey(key string) ([]byte, error) {
	kv, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, nil
	}
	return kv.Value, nil
}

func (c *Configuration) DelKey(key string) error {
	_, err := c.client.KV().Delete(key, nil)
	return err
}

func init() {
	plugin.Register(Name, &Configuration{})
}
