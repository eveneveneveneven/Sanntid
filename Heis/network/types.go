package network

// Constant parameters for the network
const (
	UDP_PORT = 20011
	TCP_PORT = 30011

	MAX_ELEVATORS = 10
)

// Network structure for passing message between Master/slave
type networkMessage struct {
	Header int
	Bool bool
}

// Declerations for networkMessage struct
var (				
	NM_REQ_CONN = &networkMessage{1, false}
	NM_REQ_ACCE = &networkMessage{0, true}
	NM_REQ_DENI = &networkMessage{0, false}
)