package netstat

import (
	"../types"
)

type NetStatHandler struct {
	becomeMaster chan bool

	networkStatus  *types.NetworkMessage
	networkUpdates *types.NetworkMessage

	netNewNetMsg    <-chan *types.NetworkMessage
	netUpdateNetMsg chan<- *types.NetworkMessage

	elevNewElevStat <-chan *types.ElevStat
	elevNewOrder    <-chan int
	elevOrderDone   <-chan int

	orderHandlerNotify chan<- *types.NetworkMessage
	orderHandlerUpdate <-chan *types.NetworkMessage
}

func NewNetStatHandler(becomeMaster chan bool,
	netNewMsg, netUpdMsg chan *types.NetworkMessage) *NetStatHandler {
	return &NetStatHandler{
		becomeMaster: becomeMaster,

		networkStatus:  nil,
		networkUpdates: new(chan *types.NetworkMessage),

		netNewNetMsg:    netNewMsg,
		netUpdateNetMsg: netUpdMsg,

		elevNewElevStat: make(chan *types.ElevStat),
		elevNewOrder:    make(chan int),
		elevOrderDone:   make(chan int),

		orderHandlerNotify: make(chan *types.NetworkMessage),
		orderHandlerUpdate: make(chan *types.NetworkMessage),
	}
}

func (ns *NetStatHandler) Run() {
	// Slave loop
slaveloop:
	for {
		// Prioritize checking the master channel
		if _, ok := <-ns.becomeMaster; !ok {
			break
		}

		select {
		case _, ok := <-ns.becomeMaster:
			if !ok {
				break slaveloop
			}
		case newMsg := <-ns.netNewNetMsg:
			ns.networkStatus = newMsg
			ns.orderHandlerNotify <- newMsg
			ns.netUpdateNetMsg <- ns.networkUpdates
			types.Clone(ns.networkUpdates, ns.networkStatus)
		case newElevStat := <-ns.elevNewElevStat:
			ns.networkUpdates.Statuses[ns.networkUpdates.Id] = *newElevStat
			ns.netUpdateNetMsg <- ns.networkUpdates
		case newOrder := <-ns.elevNewOrder:

		case orderDone := <-ns.elevOrderDone:

		}
	}

	// Master loop
	for {
		select {
		case newMsg := <-ns.netNewNetMsg:
			parseAndAddMessage(newMsg)
		case newElevStat := <-ns.elevNewElevStat:

		case newOrder := <-ns.elevNewOrder:

		case orderDone := <-ns.elevOrderDone:

		}
	}
}

func parseAndAddMessage(msg *types.NetworkMessage) {

}
