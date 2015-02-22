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

type Heis_button_press int

const (
	BUTTON_CALL_UP Heis_button_press = iota
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

func Heis_get_button(button Heis_button_press, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func Heis_stop() bool {
	return int(C.elev_get_stop_signal()) != 0
}

func Heis_obstruction() bool {
	return int(C.elev_get_obstruction_signal()) != 0
}

func Heis_set_floor_indicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func Heis_set_button_lamp(button Heis_button_press, floor int, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func Heis_set_stop_lamp(value int) {
	C.elev_set_stop_lamp(C.int(value))
}

func Heis_set_door_open_lamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}