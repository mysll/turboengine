package main

import (
	"flag"
	"os"
	"turboengine/apps/login/login"
	"turboengine/common/log"
	"turboengine/core/service"
)

var (
	config = flag.String("c", "./conf/login.toml", "config path")
)

func main() {
	flag.Parse()
	log.Init(nil)
	defer log.Close()

	log.Info("pid:", os.Getpid())
	cfg := new(service.Config)
	if err := cfg.LoadFromToml(*config); err != nil {
		panic(err)
	}
	login := service.New(new(login.Login), cfg)
	login.Start()
	login.Await()
}
