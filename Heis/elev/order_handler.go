package elev

import (
	"../types"
	"fmt"
)

type OrderHandler struct {
	currObj     *types.Order
	lastObj     *types.Order
	currNetwork *types.NetworkMessage

	netstatCurrentNetwork <-chan *types.NetworkMessage
	netstatUpdOrders      chan *types.NetworkMessage

	elevGiveNewObj chan<- *types.Order
	elevOrderUpd   <-chan *types.Order
}

func NewOrderHandler(netstatCurrNet, netstatUpdOrders chan *types.NetworkMessage,
	elevNewObj, elevOrderUpd chan *types.Order) *OrderHandler {
	return &OrderHandler{
		currObj: nil,
		lastObj: nil,

		currNetwork: types.NewNetworkMessage(),

		netstatCurrentNetwork: netstatCurrNet,
		netstatUpdOrders:      netstatUpdOrders,

		elevGiveNewObj: elevNewObj,
		elevOrderUpd:   elevOrderUpd,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {
		case updatedNetwork := <-oh.netstatCurrentNetwork:
			oh.parseNewNetwork(updatedNetwork)
		case updatedOrder := <-oh.elevOrderUpd:
			oh.parseUpdatedOrder(updatedOrder)
		}
	}
}

func (oh *OrderHandler) parseNewNetwork(updNet *types.NetworkMessage) {
	types.Clone(oh.currNetwork, updNet)
	elev.ClearAllLights()
	for _, etg := range oh.currNetwork.Statuses[oh.currNetwork.Id].InternalOrders {
		if etg == -1 {
			break
		}
		elev.SetOrderLight(&types.Order{
			ButtonPress: types.BUTTON_INTERNAL,
			Floor:       etg,
			Completed:   false,
		})
	}
	for order := range oh.currNetwork.Orders {
		elev.SetOrderLight(&order)
	}
	oh.currObj = costFunction(oh.currNetwork)
	if oh.currObj == nil {
		oh.lastObj = nil
		return
	}
	fmt.Printf("@@@> CurrObj : %+v\n@@@> LastObj : %+v\n", oh.currObj, oh.lastObj)
	if oh.lastObj == nil || oh.currObj.Floor != oh.lastObj.Floor {
		fmt.Println("@@@> TRYING TO SEND")
		oh.elevGiveNewObj <- oh.currObj
		fmt.Println("@@@> YOLO")
		oh.lastObj = oh.currObj
	}
}

func (oh *OrderHandler) parseUpdatedOrder(updOrder *types.Order) {
	fmt.Println("===> ORDER HANDLER")
	if updOrder.Completed {
		oh.deleteOrder(updOrder)
	} else if oh.newThenAdd(updOrder) {
	} else {
		return
	}
	select {
	case oh.netstatUpdOrders <- oh.currNetwork:
	case <-oh.netstatUpdOrders:
		fmt.Println("===> FLUSING")
		oh.netstatUpdOrders <- oh.currNetwork
	}
}

func (oh *OrderHandler) deleteOrder(order *types.Order) {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		internal := oh.currNetwork.Statuses[0].InternalOrders
		internal = append(internal[1:], -1)
		oh.currNetwork.Statuses[0].InternalOrders = internal
	} else {
		order.Completed = false
		delete(oh.currNetwork.Orders, *order)
	}
	if oh.lastObj != nil {
		oh.lastObj = nil
	}
}

func (oh *OrderHandler) func_name() {

}

func (oh *OrderHandler) newThenAdd(order *types.Order) bool {
	if order.ButtonPress == types.BUTTON_INTERNAL {
		newEtg := order.Floor
		internal := oh.currNetwork.Statuses[0].InternalOrders
		for i, etg := range internal {
			if newEtg == etg {
				break
			} else if etg == -1 {
				internal[i] = newEtg
				return true
			}
		}
	} else {
		if _, ok := oh.currNetwork.Orders[*order]; !ok {
			oh.currNetwork.Orders[*order] = struct{}{}
			return true
		}
	}
	return false
}
