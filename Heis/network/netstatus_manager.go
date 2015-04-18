package network

import (
	"../types"
)

type netStatManager struct {
	netstat *types.NetworkMessage

	newMsg chan *types.NetworkMessage
	update chan *types.NetworkMessage
}

func newNetStatManager(newMsg, update chan *types.NetworkMessage) *netStatManager {
	return &netStatManager{
		netstat: types.NewNetworkMessage(),

		newMsg: newMsg,
		update: update,
	}
}

func (ns *netStatManager) run() {
	for {
		select {
		case newMsg := <-ns.newMsg:
			ns.parseNewMsg(newMsg)
		case <-ns.update:
			ns.sendUpdate()
		}
	}
}

func (ns *netStatManager) parseNewMsg(msg *types.NetworkMessage) {
	id := msg.Id
	ns.netstat.Statuses[id] = msg.Statuses[id]
	for order := range msg.Orders {
		if order.Completed {
			delete(ns.netstat.Orders, order)
		} else if _, ok := ns.netstat.Orders[order]; !ok {
			ns.netstat.Orders[order] = struct{}{}
		}
	}
}

func (ns *netStatManager) sendUpdate() {
	nm := types.NewNetworkMessage()
	types.Clone(nm, ns.netstat)
	ns.update <- nm
	masterStat := ns.netstat.Statuses[0]
	ns.netstat.Statuses = make(map[int]types.ElevStat)
	ns.netstat.Statuses[0] = masterStat
}
