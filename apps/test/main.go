package main

import (
	"fmt"
	"unsafe"
)

func main() {
	i := 1
	fmt.Println(&i)
	j := (uintptr)(unsafe.Pointer(&i))

	fmt.Printf("%X", j)
}
