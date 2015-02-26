package network

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	MAXUINT = ^uint(0)
)

// Network structure for passing message between Master/slave
type networkMessage struct {
	Header int
	Bool   bool
	ID     int
}
