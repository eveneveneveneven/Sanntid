package elev

import (
	"../driver"
	"../types"
	"fmt"
	"time"
)

const (
	SPEED         = 150
	BUFFER_ORDERS = 10
	DOOR_TIMER    = 500 // milliseconds
)

type Elevator struct {
	state *types.ElevStat
	obj   *types.Order

	dirLn   chan int
	floorLn chan int

	newObj  <-chan *types.Order
	objDone chan bool

	notifyOrder chan<- *types.Order
	newElevStat chan<- *types.ElevStat
}

func NewElevator(newObj chan *types.Order,
	order chan *types.Order, elevStat chan *types.ElevStat) *Elevator {
	driver.Heis_init()
	el := &Elevator{
		state: types.NewElevStat(),
		obj:   nil,

		dirLn:   make(chan int),
		floorLn: make(chan int),

		newObj:  newObj,
		objDone: make(chan bool),

		notifyOrder: order,
		newElevStat: elevStat,
	}
	go floorIndicator()
	ClearAllLights()
	el.elevInit()
	go buttonListener(el.notifyOrder)
	go el.floorListener()
	return el
}

func (el *Elevator) Run() {
	fmt.Println("Elevator Run")
	var objQuit chan bool = nil
	for {
		select {
		case newObj := <-el.newObj:
			fmt.Printf("NEW OBJ %+v\n", newObj)
			if el.obj != nil {
				objQuit <- true
			}
			objQuit = make(chan bool)
			el.obj = newObj
			go el.goToObjective(objQuit)
		case <-el.objDone:
			fmt.Println("Obj done")
			el.openDoors()
			el.obj.Completed = true
			el.notifyOrder <- el.obj
			el.obj = nil
		case newDir := <-el.dirLn:
			fmt.Println("New direction", newDir)
			el.state.Dir = newDir
			el.newElevStat <- el.state
		case newFloor := <-el.floorLn:
			fmt.Println("New floor", newFloor)
			el.state.Floor = newFloor
			el.newElevStat <- el.state
		}
	}
}

func (el *Elevator) elevInit() {
	fmt.Println("Elev init")
	driver.Heis_set_speed(0)
	if driver.Heis_get_floor() == -1 {
		driver.Heis_set_speed(-SPEED)
		var floor int
		for {
			floor = driver.Heis_get_floor()
			if floor != -1 {
				fmt.Println("Got a floor : ", floor)
				break
			}
		}
		driver.Heis_set_speed(0)
		el.state.Floor = floor
	}
	fmt.Println("Elevator init done")
}

func (el *Elevator) floorListener() {
	var floor, currFloor int = -1, -1
	for {
		floor = driver.Heis_get_floor()
		if floor != -1 && floor != currFloor {
			el.floorLn <- floor
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
