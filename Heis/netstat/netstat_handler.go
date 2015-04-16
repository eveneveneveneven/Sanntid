package netstat

import (
	"../types"
	"fmt"
)

type NetstatHandler struct {
	becomeMaster chan bool

	networkStatus  *types.NetworkMessage
	networkUpdates *types.NetworkMessage

	nethubNewNetMsg    <-chan *types.NetworkMessage
	nethubUpdateNetMsg chan<- *types.NetworkMessage

	elevNewElevStat <-chan *types.ElevStat
	elevOrder       <-chan *types.Order

	orderhandlerNotify chan<- *types.NetworkMessage
}

func NewNetstatHandler(becomeMaster chan bool,
	netNewMsg, netUpdMsg chan *types.NetworkMessage,
	orderNotify chan *types.NetworkMessage,
	elevNewStat chan *types.ElevStat, elevOrder chan *types.Order) *NetstatHandler {
	ns := &NetstatHandler{
		becomeMaster: becomeMaster,

		networkStatus:  types.NewNetworkMessage(),
		networkUpdates: types.NewNetworkMessage(),

		nethubNewNetMsg:    netNewMsg,
		nethubUpdateNetMsg: netUpdMsg,

		elevNewElevStat: elevNewStat,
		elevOrder:       elevOrder,

		orderhandlerNotify: orderNotify,
	}
	ns.networkStatus.Statuses = append(ns.networkStatus.Statuses, *types.NewElevStat())
	ns.networkUpdates.Statuses = append(ns.networkUpdates.Statuses, *types.NewElevStat())
	return ns
}

func (ns *NetstatHandler) Run() {

	// Slave loop
slaveloop:
	for {
		select {
		case _, ok := <-ns.becomeMaster:
			if !ok {
				ns.nethubUpdateNetMsg <- ns.networkStatus
				break slaveloop
			}
		case newMsg := <-ns.nethubNewNetMsg:
			ns.slaveNewMsg(newMsg)
		case newElevStat := <-ns.elevNewElevStat:
			ns.slaveNewElevStat(newElevStat)
		case newOrder := <-ns.elevOrder:
			ns.slaveNewOrder(newOrder)
		}
	}

	// Master loop
	ns.networkStatus.Id = 0
	ns.networkStatus.Statuses[0] = ns.networkUpdates.Statuses[0]
	for {
		select {
		case newMsg := <-ns.nethubNewNetMsg:
			ns.masterNewMsg(newMsg)
		case newElevStat := <-ns.elevNewElevStat:
			ns.masterNewElevStat(newElevStat)
		case newOrder := <-ns.elevOrder:
			ns.masterNewOrder(newOrder)
		}
	}
}

func (ns *NetstatHandler) slaveNewMsg(newMsg *types.NetworkMessage) {
	ns.orderhandlerNotify <- newMsg
	types.Clone(ns.networkStatus, newMsg)
	ns.networkUpdates.Id = ns.networkStatus.Id
	ns.nethubUpdateNetMsg <- ns.networkUpdates
	for order := range ns.networkUpdates.Orders {
		delete(ns.networkUpdates.Orders, order)
	}
}

func (ns *NetstatHandler) slaveNewElevStat(newElevStat *types.ElevStat) {
	ns.networkUpdates.Statuses[0] = *newElevStat
	ns.nethubUpdateNetMsg <- ns.networkUpdates
}

func (ns *NetstatHandler) slaveNewOrder(newOrder *types.Order) {
	if _, ok := ns.networkUpdates.Orders[*newOrder]; !ok {
		fmt.Printf("Adding new order %+v\n", newOrder)
		ns.networkUpdates.Orders[*newOrder] = struct{}{}
		ns.nethubUpdateNetMsg <- ns.networkUpdates
	}
}

func (ns *NetstatHandler) masterNewMsg(newMsg *types.NetworkMessage) {
	netStat := ns.networkStatus
	if newMsg == nil {
		fmt.Printf("Netstat : %+v\n", netStat)
		ns.orderhandlerNotify <- newMsg
		stats := []types.ElevStat{netStat.Statuses[0]}
		netStat.Statuses = stats
		return
	}
	id := newMsg.Id
	numElevs := len(netStat.Statuses)
	if numElevs == id {
		netStat.Statuses = append(netStat.Statuses, newMsg.Statuses[0])
	} else if numElevs > id {
		netStat.Statuses[id] = newMsg.Statuses[0]
	} else {
		fmt.Printf(`\t\x1b[31;1mError\x1b[0m |ns.masterNewMsg| [Got id:%v,
			has only numElevs:%v],discard input\n`, id, numElevs)
	}
	for order := range newMsg.Orders {
		if _, ok := netStat.Orders[order]; !ok {
			netStat.Orders[order] = struct{}{}
		} else if order.Completed {
			delete(netStat.Orders, order)
		}
	}
	ns.nethubUpdateNetMsg <- netStat
}

func (ns *NetstatHandler) masterNewElevStat(newElevStat *types.ElevStat) {
	netStat := ns.networkStatus
	netStat.Statuses[0] = *newElevStat
	ns.nethubUpdateNetMsg <- netStat
}

func (ns *NetstatHandler) masterNewOrder(newOrder *types.Order) {
	netStat := ns.networkStatus
	if _, ok := netStat.Orders[*newOrder]; !ok {
		netStat.Orders[*newOrder] = struct{}{}
	} else if newOrder.Completed {
		delete(netStat.Orders, *newOrder)
	}
	ns.nethubUpdateNetMsg <- netStat
}
