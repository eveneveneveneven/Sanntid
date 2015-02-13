package network

type TCPHub struct {
	masterIP string
}

// Init of new TCPHub variable
func newTCPHub() *TCPHub {
	var t TCPHub

	t.masterIP = ""

	return &t
}

// Asks found Master if it can connect to the network.
// Connects itself to the network if approved,
// else shuts program off (not needed/allowed).
// Returns (isAllowed, ID, error)
func (t *TCPHub) requestConnToNetwork(masterIP string) (bool, int, error) {
	t.masterIP = masterIP
	return true, 1, nil
}