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
	objDone chan bool

	newElevstat chan *types.ElevStat

	newObj      chan *types.Order
	objComplete chan *types.Order
}

func newElevator(newElevstatCh chan *types.ElevStat,
	newObjCh, objCompleteCh chan *types.Order) *Elevator {
	driver.Heis_init()
	el := &Elevator{
		state: types.NewElevStat(),
		obj:   nil,

		dirLn:   make(chan int),
		floorLn: make(chan int),
		objDone: make(chan bool),

		newElevstat: newElevstatCh,

		newObj:      newObjCh,
		objComplete: objCompleteCh,
	}
	go floorIndicator()
	clearAllLights()
	el.elevInit()
	go el.floorListener()
	return el
}

func (el *Elevator) run() {
	fmt.Println("Start Elevator!")
	var objQuit chan bool = nil
	for {
		select {
		case obj := <-el.newObj:
			fmt.Println("elev new obj")
			if el.obj != nil {
				fmt.Println("elev quitting curr obj")
				objQuit <- true
			}
			objQuit = make(chan bool)
			el.obj = obj
			fmt.Println("elev goto obj")
			go el.goToObjective(objQuit)
		case <-el.objDone:
			el.openDoors()
			el.obj.Completed = true
			el.objComplete <- el.obj
			el.obj = nil
		case newDir := <-el.dirLn:
			el.state.Dir = newDir
			el.newElevstat <- el.state
		case newFloor := <-el.floorLn:
			el.state.Floor = newFloor
			el.newElevstat <- el.state
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
