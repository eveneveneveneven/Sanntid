package elev

import (
	"fmt"
	"sort"
	"time"

	"../backup"
	"../driver"
	"../types"
)

const (
	SPEED         = 150
	BUFFER_ORDERS = 10
	DOOR_TIMER    = 3000 // milliseconds
)

type Elevator struct {
	state *types.ElevStat
	obj   *types.Order

	dirLn   chan int
	floorLn chan int
	objDone chan bool
	stop    chan bool

	newElevstat chan types.ElevStat

	newObj      chan types.Order
	objComplete chan types.Order
}

func newElevator(newElevstatCh chan types.ElevStat,
	newObjCh, objCompleteCh chan types.Order) *Elevator {
	driver.Heis_init()
	el := &Elevator{
		state: types.NewElevStat(),
		obj:   nil,

		dirLn:   make(chan int),
		floorLn: make(chan int),
		objDone: make(chan bool),
		stop:    make(chan bool),

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
	for {
		select {
		case obj := <-el.newObj:
			if el.obj != nil {
				close(el.stop)
				el.stop = make(chan bool)
			}
			el.obj = &obj
			go el.goToObjective(el.stop, obj, el.state.Floor)
		case <-el.objDone:
			fmt.Println("sending objdone")
			el.objComplete <- *el.obj
			fmt.Println("done objdone")
			el.state.Floor = el.obj.Floor
			close(el.stop)
			el.stop = make(chan bool)
			go el.goDirection(types.STOP)
			openDoors()
			el.obj = nil
		case newDir := <-el.dirLn:
			el.state.Dir = newDir
			select {
			case el.newElevstat <- *el.state:
			case <-el.newElevstat:
				el.newElevstat <- *el.state
			}
		case newFloor := <-el.floorLn:
			el.state.Floor = newFloor
			select {
			case el.newElevstat <- *el.state:
			case <-el.newElevstat:
				el.newElevstat <- *el.state
			}
		}
	}
}

func (el *Elevator) elevInit() {
	fmt.Println("Elevator init")
	driver.Heis_set_speed(0)
	time.Sleep(100 * time.Millisecond)
	if driver.Heis_get_floor() == -1 {
		driver.Heis_set_speed(-SPEED)
		var floor int
		for {
			floor = driver.Heis_get_floor()
			if floor != -1 && floor < M_FLOORS {
				break
			}
		}
		driver.Heis_set_speed(0)
		el.state.Floor = floor
		el.state.Dir = types.STOP
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

func (el *Elevator) goToObjective(stop chan bool, obj types.Order, currFloor int) {
	dest := obj.Floor
	diff := dest - currFloor
	if diff > 0 {
		el.goDirection(types.UP)
	} else if diff < 0 {
		el.goDirection(types.DOWN)
	} else {
		el.objDone <- true
		return
	}

	for driver.Heis_get_floor() != dest {
		select {
		case _, ok := <-stop:
			if !ok {
				return
			}
		default:
		}
	}
	el.objDone <- true
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
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |el.goDirection| [Direction recieved is not valid, got the following %v], exit program\n", dir)
	}
}

func openDoors() {
	driver.Heis_set_door_open_lamp(1)
	time.Sleep(DOOR_TIMER * time.Millisecond)
	driver.Heis_set_door_open_lamp(0)
}

func processInternalBackup() {
	internal := sort.IntSlice(backup.ReadInternalBackup())
	sort.Sort(internal)
	clearAllLights()
	currFloor := driver.Heis_get_floor()
	var orders []types.Order = nil
	for i := len(internal) - 1; i >= 0; i-- {
		etg := internal[i]
		if etg != -1 {
			order := types.Order{
				ButtonPress: types.BUTTON_INTERNAL,
				Floor:       etg,
			}
			orders = append(orders, order)
			setOrderLight(&order, false)
		} else {
			break
		}
	}
	for _, order := range orders {
		diff := order.Floor - currFloor
		if diff > 0 {
			driver.Heis_set_speed(SPEED)
		} else if diff < 0 {
			driver.Heis_set_speed(-SPEED)
		}
		for currFloor != order.Floor {
			currFloor = driver.Heis_get_floor()
		}
		driver.Heis_set_speed(0)
		setOrderLight(&order, true)
		openDoors()
	}
	backup.WriteInternalBackup([]int{-1, -1, -1, -1})
}
