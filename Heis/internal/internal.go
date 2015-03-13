package internal

import (
	//"../cost"
	"../lights"
	"../buttons"
	. "../driver"
	. "fmt"
	"time"
)

var speed int
var dir int
var last_floor int
var current_order int
var queue = []int {-1, -1, -1, -1}
var costs int

func open_doors() {
	Heis_set_door_open_lamp(1)
	time.Sleep(1000 * time.Millisecond)
	Heis_set_door_open_lamp(0)
	return
}
func Init_orders() [][]int {
	orders := make([][]int, 5)
	for i := 0; i < 5; i++ {
		orders[i] = make([]int, 4)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			orders[i][j] = 0
		}
	}
	Print("Orders:")
	Println("")
	Print_orders(orders)
	return orders
}

func Print_orders(orders [][]int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			Print(orders[i][j])
			Print("     ")
		}
		println("")
	}
	println("")
}
func remove_from_queue(queue []int) []int {

	for i := 0; i < 3; i++ {
		queue[i] = queue[i+1]
	}
	queue[3] = -1
	return queue
}
func remove_from_orders(orders [][]int, current_order int) [][]int {
	if current_order != -1 {
		for i := 0; i < 3; i++ {
			orders[current_order][i] = 0
		}
	}
	lights.Set_external_lights(orders)
	lights.Set_internal_lights(orders)
	return orders
}
func Send_to_floor(orders [][]int) ([][]int) {
	for i := 0; i < 4; i++ {
		get_queue(orders)
		current_order = queue[0]
		current_floor := Heis_get_floor()
		if current_order != -1 && current_order != current_floor{
			Println("Order received: Floor ", current_order+1)
		}
		if current_order == -1 {
			Heis_set_speed(0)
			last_floor = current_floor
			dir=0
		}
		if current_floor == -1 {
			stop_all()
		}
		if Heis_get_floor() != -1 && current_order != -1 {
			Print("Current floor: ")
			Println(current_floor + 1)
			Print("Going to: Floor ")
			Println(current_order + 1)
			Println(" ")
		}
		if current_order == current_floor && current_order != -1 {
			open_doors()
			queue = remove_from_queue(queue)
			orders = remove_from_orders(orders, current_order)
			current_order = -1

		}
		if current_order > current_floor && current_order != -1 {
			Heis_set_speed(speed)
			dir=1
			for {
				if Heis_get_floor() == current_order {
					Heis_set_speed(0)
					last_floor = current_order
					dir=0
					Print("Arrived at floor: ")
					Println(current_order + 1)
					Println(" ")
					queue = remove_from_queue(queue)
					orders = remove_from_orders(orders, current_order)
					open_doors()
					current_order = -1
					return orders
				}
			}
		}
		if current_order < current_floor && current_order != -1 {
			Heis_set_speed(-speed)
			dir=-1
			for {
				if Heis_get_floor() == current_order {
					Heis_set_speed(0)
					last_floor = current_order
					dir=0
					Print("Arrived at floor: ")
					Println(current_order + 1)
					Println(" ")
					queue = remove_from_queue(queue)
					orders = remove_from_orders(orders, current_order)
					open_doors()
					current_order = -1
					return orders
				}
			}
		}
	}
	return orders
}
func get_queue(orders [][]int) {
	for i := 0; i < 4; i++ {
		queue[i] = -1
	}
	k := 0
	for i := 0; i < 4; i++ {

		if orders[i][0] == 1 || orders[i][1] == 1 || orders[i][2] == 1 {
			queue[k] = i
			k++
		}
	}
}
//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator() {
	Println("executing floor indicator!")
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
func To_nearest_floor() {
	for {

		Heis_set_speed(-speed)
		dir=-1
		if Heis_get_floor() == 0 {
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
func get_stop() {
	for {
		if Heis_stop() {
			Heis_set_stop_lamp(1)
			stop_all()
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func stop_all() {
	Heis_set_speed(0)
	dir=0
	time.Sleep(time.Millisecond * 1000)
	Heis_set_stop_lamp(0)
	To_nearest_floor()
	Heis_init()
}

func Internal() {
	// Init
	speed = 150
	orders := Init_orders()
	Heis_init()
	current_order=-1
	//Heis_set_speed(0)
	dir =0
	To_nearest_floor()
	Heis_set_stop_lamp(0)

	go get_stop()
	go Floor_indicator()
	for {
		orders = button_listener.Get_orders(orders, current_order)
		Send_to_floor(orders)
	}

	select {}
}
