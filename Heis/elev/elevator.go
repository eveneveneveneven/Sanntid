package elev

import (
	"../driver"
	"../types"
	"fmt"
	"times"
)

const (
	SPEED         = 150
	BUFFER_ORDERS = 10
	DOOR_TIMER    = 500 // milliseconds
)

type Elevator struct {
	state *types.ElevStat
	obj   *types.Order

	orderLn chan *types.Order
	dirLn   chan int
	floorLn chan int

	newObj  <-chan *types.Order
	objDone chan *types.Order

	notifyOrder chan<- *types.Order
	newElevStat chan<- *types.ElevStat
}

func NewElevator(newObj chan *types.Order,
	order chan *types.Order, elevStat chan *types.ElevStat) *Elevator {
	driver.Heis_init()
	el := &Elevator{
		state: types.NewElevStat,
		obj:   nil,

		orderLn: make(chan *types.Order, BUFFER_ORDERS),
		dirLn:   make(chan int),
		floorLn: make(chan int),

		newObj:  newObj,
		objDone: make(chan *types.Order),

		notifyOrder: order,
		newElevStat: elevStat,
	}
	go floorIndicator()
	el.elevInit()
	go buttonListener(el.orderLn)
	go el.floorListener()
	return el
}

func (el *Elevator) Run() {
	var objQuit chan bool = nil
	for {
		select {
		case newObj := <-el.newObj:
			if objQuit != nil {
				objQuit <- true
			}
			objQuit = make(chan bool)
			el.obj = newObj
			go el.goToObjective(objQuit)
		case objDone := <-el.objDone:
			el.openDoors()
			el.notifyOrder <- objDone
			el.obj = nil
		case dir := <-el.newDir:
			el.state.Dir = dir
			el.newElevStat <- el.state
		case floor := <-el.newFloor:
			el.state.Floor = floor
			el.newElevStat <- el.state
		}
	}
}

func (el *Elevator) elevInit() {
	el.goDirection(types.DOWN)
	for driver.Heis_get_floor() == -1 {
	}
	el.goDirection(types.STOP)
}

func (el *Elevator) floorListener() {
	var floor, currFloor int
	for {
		floor = driver.Heis_get_floor()
		if floor != -1 && floor != currFloor {
			el.newFloor <- floor
			currFloor = floor
		}
	}
}

func (el *Elevator) goToObjective(objQuit chan bool) {
	dest := el.obj.Floor
	diff := dest - el.state.Floor
	if diff > 0 {
		el.goDirection(types.UP)
	} else if diff < 0 {
		el.goDirection(types.DOWN)
	} else {
		el.objDone <- true
		return
	}
	stop := make(chan bool)
	go func(stop chan bool) {
		for driver.Heis_get_floor() != dest {
		}
		stop <- true
	}(stop)
	select {
	case <-objQuit:
	case <-stop:
		el.goDirection(types.STOP)
		el.objDone <- true
	}
}

func (el *Elevator) goDirection(dir int) {
	switch dir {
	case types.UP:
		driver.Heis_set_speed(SPEED)
		el.dirLn <- types.UP
	case types.DOWN:
		driver.Heis_set_speed(-SPEED)
		el.dirLn <- types.DOWN
	case types.STOP:
		driver.Heis_set_speed(0)
		el.dirLn <- types.STOP
	default:
		fmt.Printf(`\t\x1b[31;1mError\x1b[0m |el.goDirection| [Direction recieved is not valid,
			 got the following %v], exit program\n`, dir)
	}
}

func (el *Elevator) openDoors() {
	driver.Heis_set_door_open_lamp(1)
	time.Sleep(DOOR_TIMER * time.Millisecond)
	driver.Heis_set_door_open_lamp(0)
}
