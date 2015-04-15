package elev

import (
	"../buttons"
	"../cost"
	. "../driver"
	"../send_elev"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

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
