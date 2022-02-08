package cache

import (
	"fmt"
	"sync"
)

var (
	defaultItemSize = 128
	debug           = false
)

type entry struct {
	key     string
	value   any
	hashkey uint64
	next    *entry
}

type shard struct {
	lock    sync.RWMutex
	hashmap map[uint64]uint32
	items   []*entry
	count   int
	caps    int
	tail    int
	size    int
}

func newShard(size int) *shard {
	s := &shard{}
	if size <= 0 {
		size = defaultItemSize
	}
	s.size = size
	s.hashmap = make(map[uint64]uint32, size)
	s.items = make([]*entry, 0, size)
	return s
}

func (s *shard) saveItem(index uint32, e *entry) {
	if index > 0xFFFFFFFF {
		panic("index exceed")
	}

	s.caps++
	s.count++
	len := len(s.items)
	if index == uint32(len) {
		s.items = s.items[:len+1]
	}
	if index >= uint32(cap(s.items)) {
		panic("index exceed")
	}

	s.items[index] = e
	s.hashmap[e.hashkey] = uint32(index)
}

func (s *shard) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = s.items[:0]
	s.hashmap = make(map[uint64]uint32, s.size)
	s.caps = 0
	s.count = 0
	s.tail = 0
}

func (s *shard) Set(key string, hashkey uint64, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	e := &entry{
		key:     key,
		value:   value,
		hashkey: hashkey,
	}

	if index, ok := s.hashmap[hashkey]; ok {
		old := s.items[index]
		for old != nil {
			if old.key == key { // exist
				old.value = value // update
				if debug {
					fmt.Println("update index ", index)
				}
				return
			}
			if old.next == nil {
				break
			}
			old = old.next
		}
		old.next = e
		s.count++
		if debug {
			fmt.Println("link at ", index)
		}
		return
	}

	if s.tail < cap(s.items) {
		s.saveItem(uint32(s.tail), e)
		if debug {
			fmt.Println("push back at ", s.tail)
		}
		s.tail++
		return
	}

	caps := cap(s.items)
	if s.caps <= (caps >> 1) { // half empty
		s.trim(0, caps-1)
		s.tail = len(s.items)
		if debug {
			fmt.Println("trim")
			s.output()
		}
		s.saveItem(uint32(s.tail), e)
		if debug {
			fmt.Println("push back at ", s.tail)
		}
		s.tail++
		return
	}

	if s.caps < cap(s.items) {
		for k, i := range s.items {
			if i == nil {
				s.saveItem(uint32(k), e)
				if debug {
					fmt.Println("insert at ", k)
				}
				return
			}
		}
	}

	// full
	index := cap(s.items)
	s.items = append(s.items, nil) // expand
	s.tail = index + 1
	s.saveItem(uint32(index), e)
	if debug {
		fmt.Println("append back ", index)
	}
}

func (s *shard) trim(left int, right int) {
	for ; left < right; left++ {
		if s.items[left] == nil {
			break
		}
	}
	if left >= right {
		return
	}

	for ; right > left; right-- {
		if s.items[right] != nil {
			break
		}
	}
	if right <= left {
		return
	}

	hashkey := s.items[right].hashkey
	s.items[left], s.items[right] = s.items[right], s.items[left] // exchange
	s.hashmap[hashkey] = uint32(left)                             // reindex
	s.items = s.items[:right]                                     // shrink
	s.trim(left+1, right-1)
}

func (s *shard) Get(key string, hashkey uint64) any {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if index, ok := s.hashmap[hashkey]; ok {
		obj := s.items[index]
		for obj != nil {
			if obj.key == key {
				return obj.value
			}
			obj = obj.next
		}
	}
	return nil
}

func (s *shard) Del(key string, hashkey uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if index, ok := s.hashmap[hashkey]; ok {
		obj := s.items[index]
		var prev *entry
		for obj != nil {
			if obj.key == key {
				if prev == nil { // first
					if obj.next == nil {
						s.items[index] = nil
						delete(s.hashmap, hashkey)
						s.caps--
					} else {
						s.items[index] = obj.next
					}

					s.count--
					return
				}
				prev.next = obj.next
				s.count--
				return
			}
			prev = obj
			obj = obj.next
		}
	}
}

func (s *shard) output() {
	fmt.Println("**************output***************")
	fmt.Println(" tail:", s.tail, "caps:", s.caps, "count:", s.count)
	for k, v := range s.items {
		if v == nil {
			continue
		}
		fmt.Println(" index:", k, "item:\t", *v)
		obj := v.next
		for obj != nil {
			fmt.Println("\t\t", *obj)
			obj = obj.next
		}
	}
	fmt.Println("===================================")
}
