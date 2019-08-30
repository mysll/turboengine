package main

import (
	"flag"
	"os"
	"runtime/debug"
	"turboengine/apps/gate/gate"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/gate.toml", "config path")
)

func main() {
	debug.SetTraceback("single")
	flag.Parse()
	log.Init(nil)
	defer log.Close()

	log.Info("pid:", os.Getpid())
	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}

	srv := service.New(new(gate.Gate), cfg)
	if err := srv.Start(); err != nil {
		panic(err)
	}
	srv.Await()
}
