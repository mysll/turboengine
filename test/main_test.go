package main

import (
	"testing"
)

func BenchmarkChan(b *testing.B) {
	ch := make(chan int, 1000)
	go func() {
		for {
			ch <- 0
		}
	}()

	for i := 0; i < b.N; i++ {
		<-ch
	}
}

func BenchmarkWriteChan(b *testing.B) {
	ch := make(chan int, 1000)
	go func() {
		for {
			<-ch
		}
	}()

	for i := 0; i < b.N; i++ {
		ch <- i
	}
}

func BenchmarkWriteSelect(b *testing.B) {
	ch := make(chan int, 100)

	for i := 0; i < b.N; i++ {
		select {
		case <-ch:
		default:
		}
	}
}
