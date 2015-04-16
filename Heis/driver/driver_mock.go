package driver

import (
	"../types"
	"fmt"
	"math/rand"
	"time"
)

var r *rand.Rand
var ti0, ti1 time.Time

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	ti0 = time.Now()
}

func Heis_init() bool {
	return true
}

var t0, t1 time.Time
var curr_speed int

func Heis_set_speed(speed int) {
	fmt.Printf("Elevator speed set to %v\n", speed)
	t0 = time.Now()
	curr_speed = speed
}

var curr_floor int = 3

func Heis_get_floor() int {
	t1 = time.Now()
	if t1.Sub(t0).Seconds() > 1.0 {
		if curr_speed > 0 {
			curr_floor++
		} else if curr_speed < 0 {
			curr_floor--
		}
		t0 = time.Now()
		return curr_floor
	}
	if curr_speed != 0 {
		return -1
	} else {
		return curr_floor
	}
}

func Heis_get_button(button int, floor int) bool {
	ti1 = time.Now()
	if ti1.Sub(ti0).Seconds() > 2.0 {
		ti0 = time.Now()
		return true
	}
	return false
}

var last_floor int = -1

func Heis_set_floor_indicator(floor int) {
	if floor != last_floor {
		fmt.Printf("Elevator is now on floor %v\n", floor)
		last_floor = floor
	}
}

func Heis_set_button_lamp(button int, floor int, value int) {
	fmt.Print("Elevator button lamp ")
	switch button {
	case types.BUTTON_CALL_UP:
		fmt.Print("UP")
	case types.BUTTON_CALL_DOWN:
		fmt.Print("DOWN")
	case types.BUTTON_INTERNAL:
		fmt.Print("INTERNAL")
	}
	fmt.Print(" for floor", floor, "is turned")
	if value == 0 {
		fmt.Println(" off.")
	} else {
		fmt.Println(" on.")
	}
}

func Heis_set_door_open_lamp(value int) {
	if value == 0 {
		fmt.Println("Door opens.")
	} else {
		fmt.Println("Door closes.")
	}
}
