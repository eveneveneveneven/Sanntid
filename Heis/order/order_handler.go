package order

import (
	"../types"
	"../elev"
	"fmt"
)

type OrderHandler struct {
	currObj *types.Order
	lastObj *types.Order
	currNetwork *types.NetworkMessage

	netstatCurrentNetwork <-chan *types.NetworkMessage

	elevGiveNewObj chan<- *types.Order
}

func NewOrderHandler(netstatCurrNet chan *types.NetworkMessage,
	elevNewObj chan *types.Order) *OrderHandler {
	return &OrderHandler{
		currObj: nil,
		lastObj: new(types.Order),

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
	elev.ClearAllLights()
	for _, etg := range updNet.Statuses[updNet.Id].InternalOrders {
		if etg == -1 {
			break
		}
		elev.SetOrderLight(&types.Order{
			ButtonPress: types.BUTTON_INTERNAL,
			Floor: etg,
			Completed: false,
		})
	}
	for order := range updNet.Orders {
		elev.SetOrderLight(&order)
	}
	oh.currObj = costFunction(updNet)
	if oh.currObj == nil {
		return
	}
	fmt.Printf("CurrObj : %+v\n", oh.currObj)
	if *oh.currObj != *oh.lastObj {
		oh.elevGiveNewObj <- oh.currObj
		oh.lastObj = oh.currObj
	}
}
