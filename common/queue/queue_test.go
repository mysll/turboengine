package queue

import (
	"sync"
	"testing"
)

func BenchmarkQueue(b *testing.B) {
	q := NewQueue()

	var wg sync.WaitGroup
	wg.Add(1)
	i := 0
	go func() {
		for {
			q.GetOne()
			i++
			if i >= b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		q.Put(i)
	}
	wg.Wait()
}

func BenchmarkQueueAll(b *testing.B) {
	q := NewQueue()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			d, ok := q.Get()
			if !ok {
				continue
			}
			if d[len(d)-1].(int) >= b.N-1 {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		q.Put(i)
	}
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan any, 1000)
	var wg sync.WaitGroup
	wg.Add(1)
	i := 0

	go func() {
		for {
			<-ch
			i++
			if i == b.N {
				wg.Done()
				break
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		ch <- `a`
	}

	wg.Wait()
}

func TestQueue(t *testing.T) {
	q := NewQueue()
	q.Put(0, 1, 2, 3, 4, 5, 6, 7)
	data, ok := q.Get()
	if !ok {
		t.Fatal("failed get")
	}
	for k := range data {
		if data[k].(int) != k {
			t.Fatal("failed value")
		}
	}
	q.Put(0, 1, 2, 3)
	q.GetOne()
	q.GetOne()
	q.GetOne()
	q.Put(4, 5, 6, 7)
	for i := 3; i < 8; i++ {
		val, ok := q.GetOne()
		if !ok || val.(int) != i {
			t.Fatal("failed get")
		}
	}

}
