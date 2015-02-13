package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "C/channels.h"
#include "C/io.h"
#include "C/io.c"
#include "C/elev.h"
#include "C/elev.c"
*/
import "C"

import (
	"fmt"
)

func Io_init() {
	if C.io_init() == 1 {
		fmt.Println("Driver did work!")
	} else {
		fmt.Println("Driver did not work!")
	}
}