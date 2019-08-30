package main

import (
	"flag"
	"turboengine/apps/management/monitor"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/main.toml", "config path")
)

func main() {
	flag.Parse()

	log.Init(&log.Config{
		Family: "default",
		Stdout: false,
		Dir:    "./logs",
		V:      0,
	})
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}
	srv := service.New(new(monitor.Monitor), cfg)
	if err := srv.Start(); err != nil {
		panic(err)
	}
	srv.Await()
}
