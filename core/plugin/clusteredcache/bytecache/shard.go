package bytecache

import (
	"fmt"
	"sync"
)

var (
	defaultItemSize = 128
	debug           = true
)

type entry struct {
	key     string
	index   int
	hashkey uint64
	next    *entry
}

type shard struct {
	lock    sync.RWMutex
	hashmap map[uint64]uint32
	entries *BytesQueue
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
	s.entries = NewBytesQueue(1024*1024, 1024*1024*1024, false)
	return s
}

func (s *shard) saveItem(index uint32, e *entry, data []byte) error {
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

	i, err := s.entries.Push(data)
	if err != nil {
		return err
	}
	e.index = i
	s.items[index] = e
	s.hashmap[e.hashkey] = uint32(index)
	return nil
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

func (s *shard) Set(key string, hashkey uint64, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	e := &entry{
		key:     key,
		hashkey: hashkey,
	}

	if index, ok := s.hashmap[hashkey]; ok {
		old := s.items[index]
		for old != nil {
			if old.key == key { // exist
				oldbyte, err := s.entries.Get(old.index)
				if err != nil {
					panic(err)
				}
				if len(oldbyte) == len(value) {
					copy(oldbyte, value)
				} else {
					index, err := s.entries.Push(value)
					if err != nil {
						return err
					}
					old.index = index
				}

				if debug {
					fmt.Println("update index ", index)
				}
				return nil
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
		return nil
	}

	if s.tail < cap(s.items) {
		if err := s.saveItem(uint32(s.tail), e, value); err != nil {
			return err
		}
		if debug {
			fmt.Println("push back at ", s.tail)
		}
		s.tail++
		return nil
	}

	caps := cap(s.items)
	if s.caps <= (caps >> 1) { // half empty
		s.trim(0, caps-1)
		s.tail = len(s.items)
		if debug {
			fmt.Println("trim")
			s.output()
		}
		if err := s.saveItem(uint32(s.tail), e, value); err != nil {
			return err
		}
		if debug {
			fmt.Println("push back at ", s.tail)
		}
		s.tail++
		return nil
	}

	if s.caps < cap(s.items) {
		for k, i := range s.items {
			if i == nil {
				if err := s.saveItem(uint32(k), e, value); err != nil {
					return err
				}
				if debug {
					fmt.Println("insert at ", k)
				}
				return nil
			}
		}
	}

	// full
	index := cap(s.items)
	s.items = append(s.items, nil) // expand
	s.tail = index + 1
	if err := s.saveItem(uint32(index), e, value); err != nil {
		return err
	}
	if debug {
		fmt.Println("append back ", index)
	}
	return nil
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

func (s *shard) Get(key string, hashkey uint64) []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if index, ok := s.hashmap[hashkey]; ok {
		obj := s.items[index]
		for obj != nil {
			if obj.key == key {
				b, err := s.entries.Get(obj.index)
				if err != nil {
					panic(err)
				}
				return b
			}
			obj = obj.next
		}
	}
	return nil
}

func (s *shard) Del(key string, hashkey uint64) bool {
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
					return true
				}
				prev.next = obj.next
				s.count--
				return true
			}
			prev = obj
			obj = obj.next
		}
	}
	return false
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
