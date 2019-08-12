package main

import (
	"flag"
	"os"
	"turboengine/apps/gate/gate"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/gate.toml", "config path")
)

func main() {
	//service.Capture()
	flag.Parse()
	log.Init(nil)
	defer log.Close()

	log.Info("pid:", os.Getpid())
	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}
	gate := service.New(new(gate.Gate), cfg)
	gate.Start()
	gate.Await()
}
