package types

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	MAXINT = ^int(0)

	SEND_INTERVAL       = 250 // milliseconds
	BUFFER_MSG_RECIEVED = 20
)

// Network structure for passing message between Master/slave
type NetworkMessage struct {
	Id     int
	Status string
	Orders string
}
