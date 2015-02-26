package network

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	MAXINT = ^int(0)
)

// Network structure for passing message between Master/slave
type networkMessage struct {
	Id     int
	Status string
	Orders string
}
