package types

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	BACKUP_PORT = 40011

	SEND_INTERVAL       = 100 // milliseconds
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
	InternalOrders []int
}

type NetworkMessage struct {
	Id       int
	Statuses map[int]ElevStat
	Orders   map[Order]struct{}
}
