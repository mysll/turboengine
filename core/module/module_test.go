package module_test

import (
	"fmt"
	"testing"
	"time"
	"turboengine/common/log"
	"turboengine/common/utils"
	"turboengine/core/module"
	"turboengine/core/service"
)

type EchoModule struct {
	module.Module
}

func (m *EchoModule) Name() string {
	return "echo"
}

func (m *EchoModule) OnUpdate(t *utils.Time) {
	if t.FrameCount()%t.FixedFPS() == 0 {
		fmt.Println("update", t.FrameCount(), ",", t.FixedFPS())
	}
	//fmt.Println("update")
}

type Echo struct {
	service.Service
}

func TestNew(t *testing.T) {
	log.Init(nil)
	defer log.Close()
	m := module.NewWithConfig(&EchoModule{}, module.Config{
		Name:  "echo",
		Async: true,
		FPS:   10,
	})
	e := &Echo{}
	c := &service.Config{Name: "echo"}
	s := service.New(e, c)
	s.AddModule(m)
	go s.Start()
	time.Sleep(time.Second * 10)
	s.Close()
	time.Sleep(time.Second)
}
