package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"turboengine/apps/testclient/client"

	"github.com/mysll/toolkit"
)

var failed int64

func Test() {
	c := client.NewClient()
	if !c.Connect("127.0.0.1", 2001) {
		atomic.AddInt64(&failed, 1)
		return
	}
	defer c.Close()
	if !c.Login(fmt.Sprintf("test%d", toolkit.RandRange(0, 1000)), "123") {
		atomic.AddInt64(&failed, 1)
		return
	}
	if !c.WaitLogin() {
		atomic.AddInt64(&failed, 1)
		fmt.Println("failed")
	}
}

func main() {
	wg := sync.WaitGroup{}
	st := time.Now()
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * time.Duration(toolkit.RandRange(1, 1000)))
			Test()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(st).Seconds())
	fmt.Println("failed:", failed)
}
