package elev

import (
	"fmt"

	"../types"
)

const (
	M_FLOORS = 4
)

type ElevatorHub struct {
	cleanup chan bool

	currNetwork  *types.NetworkMessage
	currElevstat *types.ElevStat
	currObj      *types.Order
	newOrders    map[types.Order]bool

	msgSend    chan *types.NetworkMessage
	msgRecieve chan *types.NetworkMessage

	newNetwork chan types.NetworkMessage
	newObj     chan types.Order

	newElevstat chan types.ElevStat
	sendElevObj chan types.Order
	objComplete chan types.Order

	buttonPress chan types.Order

	elev *Elevator
}

func NewElevatorHub(cleanupCh chan bool,
	sendCh, recieveCh chan *types.NetworkMessage) *ElevatorHub {
	eh := &ElevatorHub{
		cleanup: cleanupCh,

		currNetwork:  types.NewNetworkMessage(),
		currElevstat: types.NewElevStat(),
		currObj:      nil,
		newOrders:    make(map[types.Order]bool),

		msgSend:    sendCh,
		msgRecieve: recieveCh,

		newNetwork: make(chan types.NetworkMessage),
		newObj:     make(chan types.Order),

		newElevstat: make(chan types.ElevStat, 1),
		sendElevObj: make(chan types.Order, 1),
		objComplete: make(chan types.Order, 1),

		buttonPress: make(chan types.Order),

		elev: nil,
	}
	go newOrderHandler(eh.newNetwork, eh.newObj).run()
	go buttonListener(eh.buttonPress)
	eh.elev = newElevator(eh.newElevstat, eh.sendElevObj, eh.objComplete)
	go eh.elev.run()
	return eh
}

func (eh *ElevatorHub) Run() {
	fmt.Println("Start ElevatorHub!")
	for {
		select {
		case <-eh.cleanup:
			eh.elev.goDirection(types.STOP)
			fmt.Println("Program is quitting")
			return
		case obj := <-eh.newObj:
			fmt.Println("parseNewObj")
			eh.parseNewObj(obj)
			fmt.Println("parseNewObj done")
		case obj := <-eh.objComplete:
			fmt.Println("parseObjComplete")
			eh.parseObjComplete(obj)
			fmt.Println("parseObjComplete done")
		case netstat := <-eh.msgRecieve:
			fmt.Println("msgRecieve")
			eh.parseNewMsg(netstat)
			fmt.Println("msgRecieve done")
		case elevstat := <-eh.newElevstat:
			fmt.Println("newElevstat")
			eh.parseNewElevstat(&elevstat)
			fmt.Println("newElevstat done")
		case order := <-eh.buttonPress:
			fmt.Println("buttonPress done")
			eh.parseButtonPress(order)
			fmt.Println("buttonPress done")
		}
	}
}

func (eh *ElevatorHub) parseNewObj(obj types.Order) {
	fmt.Println("NEW OBJ ::::", obj)
	eh.currObj = &obj
	select {
	case eh.sendElevObj <- obj:
	case <-eh.sendElevObj:
		fmt.Println("flushing old obj, sending updated one.")
		eh.sendElevObj <- obj
	}
}

func (eh *ElevatorHub) parseObjComplete(obj types.Order) {
	eh.newOrders[obj] = true
	if obj.ButtonPress != types.BUTTON_INTERNAL {
		delete(eh.currNetwork.Orders, obj)
		obj.ButtonPress = types.BUTTON_INTERNAL
		eh.newOrders[obj] = true
	} else {
		eh.mergeOrder(eh.currNetwork, obj, true)
		newObj := costFunction(eh.currNetwork)
		if newObj != nil && newObj.Floor == obj.Floor {
			delete(eh.currNetwork.Orders, *newObj)
			eh.newOrders[*newObj] = true
		}
	}
}

func (eh *ElevatorHub) parseNewMsg(netstat *types.NetworkMessage) {
	fmt.Println("recieved ::", netstat)
	eh.currNetwork = netstat
	setActiveLights(netstat)
	response := types.NewNetworkMessage()
	response.Id = netstat.Id
	respElevStat := *types.NewElevStat()
	internal := netstat.Statuses[netstat.Id].InternalOrders
	if len(internal) != 0 {
		for i, v := range internal {
			respElevStat.InternalOrders[i] = v
		}
	}
	respElevStat.Dir = eh.currElevstat.Dir
	respElevStat.Floor = eh.currElevstat.Floor
	response.Statuses[response.Id] = respElevStat
	eh.removeRedundantOrders(response)
	eh.currElevstat.InternalOrders = internal
	eh.currNetwork.Statuses[eh.currNetwork.Id] = *eh.currElevstat
	eh.newNetwork <- *eh.currNetwork
	eh.msgSend <- response
	fmt.Println("currnetw ::", eh.currNetwork)
	fmt.Println("response ::", response)
}

func (eh *ElevatorHub) removeRedundantOrders(response *types.NetworkMessage) {
	for order, completed := range eh.currNetwork.Orders {
		if completed {
			delete(eh.currNetwork.Orders, order)
		}
		response.Orders[order] = completed
	}
	for newOrder, completed := range eh.newOrders {
		if completed {
			delete(eh.currNetwork.Orders, newOrder)
		}
		eh.mergeOrder(response, newOrder, completed)
		delete(eh.newOrders, newOrder)
	}
}

func (eh *ElevatorHub) mergeOrder(dst *types.NetworkMessage,
	order types.Order, value bool) {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		eh.mergeInternalOrder(dst, order, value)
	} else {
		eh.mergeExternalOrder(dst, order, value)
	}
}

func (eh *ElevatorHub) mergeInternalOrder(dst *types.NetworkMessage,
	order types.Order, value bool) {
	fmt.Println("merger order ::", order, ":: value ::", value)
	newEtg := order.Floor
	internal := dst.Statuses[dst.Id].InternalOrders
	for i, etg := range internal {
		if etg == newEtg {
			if value {
				internal = append(internal, -1)
				internal = append(internal[:i], internal[i+1:]...)
				break
			}
			break
		} else if etg == -1 {
			if !value {
				internal[i] = newEtg
			}
			break
		}
	}
	fmt.Println("and now     ::", internal)
	elevstat := dst.Statuses[dst.Id]
	elevstat.InternalOrders = internal
	var newInternal []int = nil
	for _, v := range internal {
		newInternal = append(newInternal, v)
	}
	eh.currElevstat.InternalOrders = newInternal
	dst.Statuses[dst.Id] = elevstat
}

func (eh *ElevatorHub) mergeExternalOrder(dst *types.NetworkMessage,
	order types.Order, value bool) {
	dst.Orders[order] = value
}

func (eh *ElevatorHub) parseNewElevstat(elevstat *types.ElevStat) {
	eh.currElevstat.Dir = elevstat.Dir
	eh.currElevstat.Floor = elevstat.Floor
}

func (eh *ElevatorHub) parseButtonPress(order types.Order) {
	eh.newOrders[order] = false
}

func CleanExit() {
	newElevator(nil, nil, nil)
}
