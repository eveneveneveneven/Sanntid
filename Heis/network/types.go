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

var (
	NM_REQ_CONN   = &networkMessage{ID_REQ_CONN, false, -1}
	NM_REQ_ACCEPT = &networkMessage{ID_REQ_ACCEPT, true, -1}
	NM_REQ_DENIED = &networkMessage{ID_REQ_DENIED, false, -1}
)
