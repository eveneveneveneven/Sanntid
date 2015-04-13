package internal

import (
	"../cost"
	"../buttons"
	"../send_elev"
	. "../driver"
	"fmt"
	"time"
	"bufio"
	"os"
	"strconv"
)

var speed int
var dir int
var last_floor int
var current_order int
var queue = []int {-1, -1, -1, -1}
var costs int
var ordered_floor int

//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator() {
	var floor int
	for {
		floor = Heis_get_floor()
		//Println(floor)
		if floor != -1 {
			Heis_set_floor_indicator(floor)
			//time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func Get_last_floor() {
	if Heis_get_floor() !=-1 {
		last_floor = Heis_get_floor()
	}
}

func To_nearest_floor() {
	for {

		Heis_set_speed(-speed)
		dir=-1
		if Heis_get_floor() != -1 {
			Heis_set_speed(0)
			last_floor = Heis_get_floor()
			dir=0

			return
		}
		if Heis_get_floor() != -1 {
			Heis_set_speed(0)
			last_floor = Heis_get_floor()
			dir=0
			return
		}
	}
}
func get_input() {
	reader := bufio.NewReader(os.Stdin) 
	costs := 0
	current_floor := Heis_get_floor()
	for { 
		fmt.Print("Enter int: ") 
		text, _ := reader.ReadString('\n') 
		i, err := strconv.Atoi(text[:len(text)-1])
		if err != nil || i>3 || i<0 {
			fmt.Println("Index out of range, or wrong format. Try again: ") 
		}else{
			ordered_floor = i
			current_floor = Heis_get_floor()
			if ordered_floor < current_floor{
				dir = -1
			}
			if ordered_floor > current_floor{
				dir = 1
			}
			if ordered_floor == current_floor{
				dir = 0
			}
			costs = cost.Cost_function(current_floor, ordered_floor, dir)
			fmt.Println("Cost: ", costs)
			Send_to_floor.Send_to_floor(ordered_floor)
		}
	} 
}

func Internal() {
	// Init
	//button_listener_chan := make(chan int)
	//send_elev_chan := make(chan int)
	speed = 150
	Heis_init()
	dir = 0
	To_nearest_floor()
	Heis_set_stop_lamp(0)
	go Floor_indicator()
	go get_input()
	go Get_last_floor()
	go button_listener.Get_orders()
	//<- button_listener_chan
	go Send_to_floor.Send_to_floor(ordered_floor)
	//<- send_elev_chan
	select {}
}
