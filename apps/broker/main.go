package main

import (
	"flag"
	"turboengine/apps/broker/broker"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/broker.toml", "config path")
)

func main() {
	flag.Parse()
	log.Init(nil)
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}
	srv := service.New(new(broker.Broker), cfg)
	if err := srv.Start(); err != nil {
		panic(err)
	}
	srv.Await()
}
