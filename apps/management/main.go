package main

import (
	"turboengine/apps/management/monitor"
	"turboengine/common/log"
	"turboengine/core/service"
)

func main() {
	log.Init(&log.Config{
		Family: "default",
		Stdout: false,
		Dir:    "./logs",
		V:      0,
	})
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml("./conf/main.toml"); err != nil {
		panic(err)
	}
	monitor := service.New(new(monitor.Monitor), cfg)
	monitor.Start()
	monitor.Wait()
}
