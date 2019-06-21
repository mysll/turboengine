package main

import (
	"fmt"
	"log"
	"net"
	"turboengine/common/net/rpc"
	"turboengine/common/net/smartrpc"
	"turboengine/common/net/smartrpc/math"
)

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:0") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	return l, l.Addr().String()
}

var serverAddr string

func StartServer() {
	newServer := rpc.NewServer()
	var l net.Listener
	l, serverAddr = listenTCP()
	math.SetMathProvider(newServer, "", new(math.MathService))
	go newServer.Accept(l)
}

func NewMath() *math.Math {
	c, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	math := math.NewMathConsumer(c, "")
	return math
}

func main() {
	smartrpc.Generate(new(math.Math), "math", "../common/net/smartrpc/math/")
	StartServer()
	math := NewMath()
	res, _ := math.Do(0, 0, 1, 2)
	fmt.Println(res)
	math.Print(0, 0, "hello world")
}
