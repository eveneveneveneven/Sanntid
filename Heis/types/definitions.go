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

type ElevatorStatus struct {
	Dir            Direction
	Floor          int
	InternalOrders []bool
}

type NetworkMessage struct {
	Id        int
	Statuses  []ElevatorStatus
	Orders    []int
	NewOrders []int
}
