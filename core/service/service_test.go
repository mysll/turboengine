package service

import (
	"testing"
	"turboengine/common/log"
)

type Echo struct {
	Service
}

func TestNewService(t *testing.T) {
	log.Init(nil)
	defer log.Close()
	e := &Echo{}
	c := &Config{Name: "echo"}
	s := New(e, c)
	s.Start()
}
