package network

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	MAX_ELEVATORS = 10

	NUM_DATA_LENGTH = 2
)

// Network structure for passing message between Master/slave
type networkStatusMsg struct {
	Id int
	Data string
}