package utils

import (
	"hash/fnv"
	"math"
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

const MIN = 0.000001

func IsEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < MIN
}
