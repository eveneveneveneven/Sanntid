package elev

import (
	"fmt"

	"../types"
)

type ElevatorHub struct {
	currNetwork  *types.NetworkMessage
	currElevstat *types.ElevStat
	currObj      *types.Order
	newOrders    map[*types.Order]struct{}

	becomeMaster chan bool

	msgSend    chan *types.NetworkMessage
	msgRecieve chan *types.NetworkMessage

	newNetwork chan *types.NetworkMessage
	newObj     chan *types.Order

	newElevstat chan *types.ElevStat
	sendElevObj chan *types.Order
	objComplete chan *types.Order

	buttonPress chan *types.Order
}

func NewElevatorHub(becomeMaster chan bool,
	sendCh, recieveCh chan *types.NetworkMessage) *ElevatorHub {
	eh := &ElevatorHub{
		currNetwork:  types.NewNetworkMessage(),
		currElevstat: types.NewElevStat(),
		currObj:      nil,
		newOrders:    make(map[*types.Order]struct{}),

		becomeMaster: becomeMaster,

		msgSend:    sendCh,
		msgRecieve: recieveCh,

		newNetwork: make(chan *types.NetworkMessage),
		newObj:     make(chan *types.Order),

		newElevstat: make(chan *types.ElevStat),
		sendElevObj: make(chan *types.Order),
		objComplete: make(chan *types.Order),

		buttonPress: make(chan *types.Order),
	}
	go newOrderHandler(eh.newNetwork, eh.newObj).run()
	go buttonListener(eh.buttonPress)
	go newElevator(eh.newElevstat, eh.sendElevObj, eh.objComplete).run()
	return eh
}

func (eh *ElevatorHub) Run() {
	fmt.Println("Start ElevatorHub!")
	for {
		select {
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
	eh.currObj = obj
	eh.sendElevObj <- obj
}

func (eh *ElevatorHub) parseObjComplete(obj *types.Order) {
	eh.newOrders[obj] = struct{}{}
}

func (eh *ElevatorHub) parseNewMsg(netstat *types.NetworkMessage) {
	eh.currNetwork = netstat
	eh.newNetwork <- netstat
	response := types.NewNetworkMessage()
	response.Id = netstat.Id
	response.Statuses[response.Id] = *eh.currElevstat
	for order := range netstat.Orders {
		response.Orders[order] = struct{}{}
	}
	for order := range eh.newOrders {
		if order.Completed {
			order.Completed = false
			delete(eh.currNetwork.Orders, *order)
			delete(response.Orders, *order)
			order.Completed = true
		}
		response.Orders[*order] = struct{}{}
		delete(eh.newOrders, order)
	}
	eh.msgSend <- response
}

func (eh *ElevatorHub) parseNewElevstat(elevstat *types.ElevStat) {
	eh.currElevstat = elevstat
}

func (eh *ElevatorHub) parseButtonPress(order *types.Order) {
	eh.newOrders[order] = struct{}{}
}

/*
func (eh *ElevatorHub) newThenAdd(order *types.Order) bool {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		newEtg := order.Floor
		internal := eh.currNetwork.Statuses[0].InternalOrders
		for i, etg := range internal {
			if newEtg == etg {
				break
			} else if etg == -1 {
				internal[i] = newEtg
				return true
			}
		}
	}
}

func (eh *ElevatorHub) deleteOrder(order *types.Order) {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		internal := eh.currNetwork.Statuses[0].InternalOrders
		internal = append(internal[1:], -1)
		eh.currNetwork.Statuses[0].InternalOrders = internal
	} else {
		order.Completed = false
		delete(eh.currNetwork.Orders, *order)
	}
}
*/
