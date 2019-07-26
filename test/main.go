package main

import "fmt"

func main() {

	a := make([]int, 10)
	for k := range a {
		a[k] = k
	}

	fmt.Println(a)
	a = append(a[:9], a[10:]...)
	fmt.Println(a)

}
