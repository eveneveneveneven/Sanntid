package network

type TCPHub struct {
	masterIP string
}

func newTCPHub() *TCPHub {
	var t TCPHub



	return &t
}

func (t *TCPHub) requestConnToNetwork(masterIP string) (bool, int, error) {
	t.masterIP = masterIP
	return true, 1, nil
}