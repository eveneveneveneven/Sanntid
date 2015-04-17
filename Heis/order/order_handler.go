package order

import (
	"../elev"
	"../types"
	"fmt"
)

type OrderHandler struct {
	currNetwork *types.NetworkMessage

	netstatCurrentNetwork <-chan *types.NetworkMessage

	elevGiveNewObj chan<- *types.Order
}

func NewOrderHandler(netstatCurrNet chan *types.NetworkMessage,
	elevNewObj chan *types.Order) *OrderHandler {
	return &OrderHandler{
		currNetwork: new(types.NetworkMessage),

		netstatCurrentNetwork: netstatCurrNet,

		elevGiveNewObj: elevNewObj,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {
		case updatedNetwork := <-oh.netstatCurrentNetwork:
			oh.parseNewNetwork(updatedNetwork)
		}
	}
}

func (oh *OrderHandler) parseNewNetwork(updNet *types.NetworkMessage) {
	netStat := oh.currNetwork
	if updNet != nil {
		types.Clone(netStat, updNet)
		for order := range netStat.Orders {
			elev.SetOrderLight(&order)
		}
		// Dirty hack...
		if len(netStat.Statuses) <= netStat.Id {
			return
		}
		for _, floor := range netStat.Statuses[netStat.Id].InternalOrders {
			if floor != -1 {
				elev.SetOrderLight(&types.Order{
					ButtonPress: types.BUTTON_INTERNAL,
					Floor:       floor,
					Completed:   false,
				})
			} else {
				break
			}
		}
	} else {
		fmt.Println(`\t\x1b[31;1mError\x1b[0m |ohh.parseNewNetwork| [Got updNet
			= nil], discard input\n`)
	}
}
