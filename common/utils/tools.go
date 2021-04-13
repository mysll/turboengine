package utils

import (
	"hash/fnv"
)

func Hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func Hash64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func IsPowerOfTwo(number int) bool {
	return (number & (number - 1)) == 0
}
