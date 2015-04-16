package order

import (
	"../types"
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
	if updNet != nil {
		types.Clone(oh.currNetwork, updNet)
	}
}
