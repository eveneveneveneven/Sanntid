package types

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	SEND_INTERVAL       = 250 // milliseconds
	BUFFER_MSG_RECIEVED = 20
)

// Type definitions for the elevator and networkmessage protocol
type Direction int

const (
	UP Direction = iota
	DOWN
	STOP
)

type ButtonType_t int

const (
	BUTTON_CALL_UP ButtonType_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type Button struct {
	ButtonType ButtonType_t
	Floor      int
}

type ElevStat struct {
	Dir            Direction
	Floor          int
	InternalOrders []bool
}

type NetworkMessage struct {
	Id        int
	Statuses  []ElevStat
	Orders    []int
	NewOrders []int
}
