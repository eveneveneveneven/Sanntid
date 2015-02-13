package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "C/io.h"
#include "C/channels.h"
#include "C/elev.h"
*/
import "C"

import "fmt"

func main() {
	if !bool(C.io_init()) {
		fmt.Println("Did not work!")
	} else {
		fmt.Println("Did work!")
	}
}