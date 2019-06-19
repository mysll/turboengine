package main

import (
	"flag"
	"turboengine/common/log"

	"github.com/BurntSushi/toml"
)

func main() {
	flag.Parse()
	var config log.Config
	if _, err := toml.DecodeFile("../configs/log.toml", &config); err != nil {
		panic(err)
	}

	log.Init(&config)
	defer log.Close()
	log.Info("test")
}
