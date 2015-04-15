package driver

import "../types"

/*
#cgo LDFLAGS: -lcomedi -lm
#include "C/channels.h"
#include "C/io.h"
#include "C/io.c"
#include "C/elev.h"
#include "C/elev.c"
*/
import "C"

type Elev_button_type_t int

const (
	BUTTON_CALL_UP Elev_button_type_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

func Heis_init() bool {
	return int(C.elev_init()) != 0
}

func Heis_set_speed(speed int) {
	C.elev_set_speed(C.int(speed))
}

func Heis_get_floor() int {
	return int(C.elev_get_floor_sensor_signal())
}

func Heis_get_button(button Elev_button_type_t, floor int) bool {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor))) != 0
}

func Heis_set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Heis_set_button_lamp(button Elev_button_type_t, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Heis_set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}
