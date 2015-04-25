package types

// Type definitions for the elevator and networkmessage protocol

// Directions
const (
	UP int = iota
	DOWN
	STOP
)

// Buttons
const (
	BUTTON_CALL_UP int = iota
	BUTTON_CALL_DOWN
	BUTTON_INTERNAL
)

// Orders
type Order struct {
	ButtonPress int
	Floor       int
}

// Status of the Elevator
type ElevStat struct {
	Dir            int
	Floor          int
	InternalOrders []int
}

// Status of the Network
type NetworkMessage struct {
	Id       int
	Statuses map[int]ElevStat
	Orders   map[Order]bool
}
