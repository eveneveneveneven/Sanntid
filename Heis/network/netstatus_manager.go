package network

import (
	"fmt"
	"sort"

	"../types"
)

type netStatManager struct {
	netstat *types.NetworkMessage

	newMsg chan *types.NetworkMessage
	update chan *types.NetworkMessage
	tick   chan bool
}

func newNetStatManager(newMsgCh, updateCh chan *types.NetworkMessage,
	tickCh chan bool) *netStatManager {
	ns := &netStatManager{
		netstat: types.NewNetworkMessage(),

		newMsg: newMsgCh,
		update: updateCh,
		tick:   tickCh,
	}
	ns.netstat.Id = 0
	ns.netstat.Statuses[0] = *types.NewElevStat()
	return ns
}

func (ns *netStatManager) run(currNetstat *types.NetworkMessage, resetCh chan bool) {
	fmt.Println("Start NetStatManager!")
	ns.netstat = currNetstat
	ns.netstat.Id = 0
	for {
		select {
		case _, ok := <-resetCh:
			if !ok {
				return
			}
		case newMsg := <-ns.newMsg:
			ns.parseNewMsg(newMsg)
		case <-ns.tick:
			ns.sendUpdate()
		}
	}
}

func (ns *netStatManager) parseNewMsg(msg *types.NetworkMessage) {
	var id int
	for id = range msg.Statuses {
		break
	}
	ns.netstat.Statuses[id] = msg.Statuses[id]
	for order, completed := range msg.Orders {
		if completed {
			ns.netstat.Orders[order] = true
		} else if _, ok := ns.netstat.Orders[order]; !ok {
			ns.netstat.Orders[order] = false
		}
	}
}

func (ns *netStatManager) sendUpdate() {
	nm := types.NewNetworkMessage()
	var ids sort.IntSlice = nil
	for id := range ns.netstat.Statuses {
		ids = append(ids, id)
	}
	sort.Sort(ids)
	newStatues := make(map[int]types.ElevStat)
	for i, id := range ids {
		if id-i != 0 {
			newStatues[id-i] = ns.netstat.Statuses[id]
		} else {
			newStatues[id] = ns.netstat.Statuses[id]
		}
	}
	for order, completed := range ns.netstat.Orders {
		if completed {
			delete(ns.netstat.Orders, order)
		}
	}
	ns.netstat.Statuses = newStatues
	types.DeepCopy(nm, ns.netstat)
	ns.update <- nm

	for id := range ns.netstat.Statuses {
		if id != 0 {
			delete(ns.netstat.Statuses, id)
		}
	}
}
