package main

import (
	"turboengine/apps/gate/gate"
	"turboengine/common/log"
	"turboengine/core/service"
)

func main() {
	log.Init(nil)
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml("./conf/main.toml"); err != nil {
		panic(err)
	}
	gate := service.New(new(gate.Gate), cfg)
	gate.Start()
	gate.Wait()
}
