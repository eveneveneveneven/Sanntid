package main

import (
	"fmt"
)

type P struct {
	x, y int
}

func main() {
	p := &P{1, 2}
	fmt.Printf("%+v\n", p)
}
