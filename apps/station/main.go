package main

import (
	"flag"
	"turboengine/apps/station/station"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/station.toml", "config path")
)

func main() {
	flag.Parse()
	log.Init(nil)
	defer log.Close()

	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}
	srv := service.New(new(station.Station), cfg)
	if err := srv.Start(); err != nil {
		panic(err)
	}
	srv.Await()
}
