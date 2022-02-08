package queue

import (
	"sync"
)

type item struct {
	data any
	next *item
}

type Queue struct {
	sync.Mutex
	head  *item
	tail  *item
	count int
}

func NewQueue() *Queue {
	q := &Queue{}
	return q
}

func (q *Queue) Count() int {
	var count int
	q.Lock()
	count = q.count
	q.Unlock()
	return count
}

func (q *Queue) Empty() bool {
	empty := false
	q.Lock()
	empty = q.count == 0
	q.Unlock()
	return empty
}

func (q *Queue) Put(items ...any) {
	q.Lock()
	for _, data := range items {
		q.count++
		i := &item{
			data: data,
		}
		if q.tail == nil { // empty
			q.head = i
			q.tail = q.head
			continue
		}
		q.tail.next = i
		q.tail = i
	}

	q.Unlock()
}

func (q *Queue) Get() ([]any, bool) {
	if q.count == 0 {
		return nil, false
	}

	var res []any
	q.Lock()
	if q.count > 0 {
		for q.head != nil {
			res = append(res, q.head.data)
			q.head = q.head.next
			q.count--
		}
		q.tail = nil // empty
	}
	q.Unlock()
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

func (q *Queue) GetOne() (any, bool) {
	if q.count == 0 {
		return nil, false
	}

	var res any
	q.Lock()
	if q.head != nil {
		res = q.head.data
		q.head = q.head.next
		q.count--
		if q.head == nil {
			q.tail = nil
		}
	}
	q.Unlock()
	if res == nil {
		return nil, false
	}
	return res, true
}
