package types

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	SEND_INTERVAL       = 250 // milliseconds
	BUFFER_MSG_RECIEVED = 20
)

// Type definitions for the elevator and networkmessage protocol
const (
	UP int = iota
	DOWN
	STOP
)

const (
	BUTTON_CALL_UP int = iota
	BUTTON_CALL_DOWN
	BUTTON_INTERNAL
)

type Order struct {
	ButtonPress int
	Floor       int
	Completed   bool
}

type ElevStat struct {
	Dir            int
	Floor          int
<<<<<<< HEAD
	InternalOrders []int
=======
	InternalOrders []bool
>>>>>>> 23d545ecd2d0f08b3f30f1b40de27a25060d8f4c
}

type NetworkMessage struct {
	Id       int
	Statuses []ElevStat
	Orders   map[Order]struct{}
}
