package driver

// #cgo CFLAGS: -std=c99 -g -Wall -O2 -I
// #cgo LDFLAGS: -lcomedi -g -lm
// #include "C/io.h"
// #include "C/elev.c"
// #include "C/channels.h"
import "C"

import "fmt"

func main() {
	if ok, _ := !C.io_init(); ok {
		fmt.Println("Did not work!")
	} else {
		fmt.Println("Did work!")
	}
}