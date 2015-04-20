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
		objComplete: make(chan *types.Order, 1),

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
			fmt.Println("new objective")
			eh.parseNewObj(obj)
		case obj := <-eh.objComplete:
			fmt.Println("objective complete!")
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
}

func (eh *ElevatorHub) parseNewMsg(netstat *types.NetworkMessage) {
	fmt.Println("recieved ::", netstat)
	eh.currNetwork = netstat
	eh.currNetwork.Statuses[eh.currNetwork.Id] = *eh.currElevstat
	response := types.NewNetworkMessage()
	response.Id = eh.currNetwork.Id
	response.Statuses = eh.currNetwork.Statuses
	eh.checkRedundantOrders(response)
	eh.newNetwork <- eh.currNetwork
	eh.msgSend <- response
	fmt.Println("response ::", response)
}

func (eh *ElevatorHub) checkRedundantOrders(response *types.NetworkMessage) {
	for order := range eh.currNetwork.Orders {
		response.Orders[order] = struct{}{}
	}
	for newOrder := range eh.newOrders {
		if newOrder.Completed {
			eh.removeOrder(response, newOrder)
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
