package elev

import (
	"fmt"

	"../types"
)

type ElevatorHub struct {
	cleanup chan bool

	currNetwork  *types.NetworkMessage
	currElevstat *types.ElevStat
	currObj      *types.Order
	newOrders    map[*types.Order]struct{}

	msgSend    chan *types.NetworkMessage
	msgRecieve chan *types.NetworkMessage

	newNetwork chan *types.NetworkMessage
	newObj     chan *types.Order
	reset      chan bool

	newElevstat chan *types.ElevStat
	sendElevObj chan *types.Order
	objComplete chan *types.Order

	buttonPress chan *types.Order

	elev *Elevator
}

func NewElevatorHub(cleanupCh chan bool,
	sendCh, recieveCh chan *types.NetworkMessage) *ElevatorHub {
	eh := &ElevatorHub{
		cleanup: cleanupCh,

		currNetwork:  types.NewNetworkMessage(),
		currElevstat: types.NewElevStat(),
		currObj:      nil,
		newOrders:    make(map[*types.Order]struct{}),

		msgSend:    sendCh,
		msgRecieve: recieveCh,

		newNetwork: make(chan *types.NetworkMessage),
		newObj:     make(chan *types.Order),
		reset:      make(chan bool),

		newElevstat: make(chan *types.ElevStat),
		sendElevObj: make(chan *types.Order),
		objComplete: make(chan *types.Order, 1),

		buttonPress: make(chan *types.Order),

		elev: nil,
	}
	go newOrderHandler(eh.newNetwork, eh.newObj, eh.reset).run()
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
			return
		case obj := <-eh.newObj:
			eh.parseNewObj(obj)
		case obj := <-eh.objComplete:
			eh.parseObjComplete(obj)
		case netstat := <-eh.msgRecieve:
			eh.parseNewMsg(netstat)
		case elevstat := <-eh.newElevstat:
			eh.parseNewElevstat(elevstat)
		case order := <-eh.buttonPress:
			eh.parseButtonPress(order)
		}
	}
}

func (eh *ElevatorHub) parseNewObj(obj *types.Order) {
	fmt.Println("parseNewObj")
	eh.currObj = obj
	select {
	case eh.sendElevObj <- obj:
	default:
		fmt.Println("elevator trying to communicate with elevhub")
		eh.reset <- true
	}
	fmt.Println("parseNewObj done")
}

func (eh *ElevatorHub) parseObjComplete(obj *types.Order) {
	fmt.Println("parseObjComplete")
	eh.newOrders[obj] = struct{}{}
	eh.removeOrder(eh.currNetwork, obj)
	internal := &types.Order{
		ButtonPress: types.BUTTON_INTERNAL,
		Floor:       obj.Floor,
		Completed:   true,
	}
	eh.removeOrder(eh.currNetwork, internal)
	for {
		newObj := costFunction(eh.currNetwork)
		if newObj != nil && newObj.Floor == obj.Floor {
			newObj.Completed = true
			eh.newOrders[newObj] = struct{}{}
			eh.removeOrder(eh.currNetwork, newObj)
		} else {
			break
		}
	}
	fmt.Println("parseObjComplete done")
}

func (eh *ElevatorHub) parseNewMsg(netstat *types.NetworkMessage) {
	fmt.Println("recieved ::", netstat)
	eh.currNetwork = netstat
	eh.currNetwork.Statuses[eh.currNetwork.Id] = *eh.currElevstat
	response := types.NewNetworkMessage()
	response.Id = eh.currNetwork.Id
	response.Statuses = eh.currNetwork.Statuses
	eh.removeRedundantOrders(response)
	for respOrder := range response.Orders {
		if respOrder.Completed {
			eh.removeOrder(eh.currNetwork, &respOrder)
		}
	}
	eh.newNetwork <- eh.currNetwork
	eh.msgSend <- response
	fmt.Println("response ::", response)
}

func (eh *ElevatorHub) removeRedundantOrders(response *types.NetworkMessage) {
	for order := range eh.currNetwork.Orders {
		response.Orders[order] = struct{}{}
	}
	for newOrder := range eh.newOrders {
		if newOrder.Completed {
			eh.removeOrder(response, newOrder)
			eh.removeOrder(eh.currNetwork, newOrder)
			if newOrder.ButtonPress != types.BUTTON_INTERNAL {
				eh.addOrder(response, newOrder)
			}
		} else {
			eh.addOrder(response, newOrder)
		}
		delete(eh.newOrders, newOrder)
	}
}

func (eh *ElevatorHub) addOrder(dst *types.NetworkMessage, order *types.Order) {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		newEtg := order.Floor
		for i, etg := range eh.currElevstat.InternalOrders {
			if etg == newEtg {
				break
			} else if etg == -1 {
				eh.currElevstat.InternalOrders[i] = newEtg
				break
			}
		}
		dst.Statuses[dst.Id] = *eh.currElevstat
	} else {
		dst.Orders[*order] = struct{}{}
	}
}

func (eh *ElevatorHub) removeOrder(dst *types.NetworkMessage, order *types.Order) {
	fmt.Println("removing order ::", order)
	order.Completed = false
	if order.ButtonPress == types.BUTTON_INTERNAL {
		newEtg := order.Floor
		internal := eh.currElevstat.InternalOrders
		for i, etg := range internal {
			if etg == newEtg {
				internal = append(internal, -1)
				internal = append(internal[:i], internal[i+1:]...)
				break
			} else if etg == -1 {
				break
			}
		}
		eh.currElevstat.InternalOrders = internal
		dst.Statuses[dst.Id] = *eh.currElevstat
	} else {
		delete(dst.Orders, *order)
	}
	order.Completed = true
}

func (eh *ElevatorHub) parseNewElevstat(elevstat *types.ElevStat) {
	eh.currElevstat.Dir = elevstat.Dir
	eh.currElevstat.Floor = elevstat.Floor
}

func (eh *ElevatorHub) parseButtonPress(order *types.Order) {
	eh.newOrders[order] = struct{}{}
}

func CleanExit() {
	newElevator(nil, nil, nil)
}
