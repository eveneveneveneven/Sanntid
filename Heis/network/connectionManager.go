package network

import (
	"net"
	"sync"
)

type connManager struct {
	masterIP string
	currId   uint
	idList   map[*net.TCPConn]int

	wakeRecieve chan *networkMessage
	wakeSend    chan *networkMessage
	connEnd     chan *net.TCPConn

	hubRecieve chan *networkMessage
	hubSend    chan *networkMessage

	wg *sync.WaitGroup
}

func NewConnManager(hbRec, hbSend chan *networkMessage) *connManager {
	var cm connManager

	cm.masterIP = ""
	cm.numConns = 1
	cm.conns = make(map[int]*net.TCPConn)

	cm.wakeRecieve = make(chan *networkMessage, 20) // buffer for messages recieved
	cm.wakeSend = make(chan *networkMessage)
	cm.connEnd = make(chan *net.TCPConn)

	hubRecieve = hbRec
	hubSend = hbSend

	return &cm
}

func (cm *connManager) run() {
	for {
		// prioritized channel to check
		select {
		case conn := <-connEnd:
			cm.removeConnection(conn)
			continue
		default:
		}

		select {
		case conn := <-connEnd:
			cm.removeConnection(conn)
		case recieveMsg := <-wakeRecieve:
			hubRecieve <- recieveMsg
		case sendMsg := <-cm.hubSend:
			numConns = len(cm.conns)
			wg.Add(numConns)
			for i := 0; i < numConns; i++ {
				wakeRecieve <- sendMsg
			}
			wg.Wait()
		}
	}
}

func (cm *connManager) connectToNetwork(masterIP string) (int, error) {
	cm.masterIP = masterIP
	conn, err := createConnTCP(cm.masterIP)
	if err != nil {
		return -1, err
	}
	cm.addConnection(conn)
	go createTCPHandler(conn, cm.wakeRecieve, cm.wakeSend, cm.connEnd, cm.wg)
}

func (cm *connManager) addConnection(conn *net.TCPConn) {
	cm.conns[conn] = currId
	currId++
}

func (cm *connManager) removeConnection(conn *net.TCPConn) {
	if removeId, ok := cm.conns[conn]; ok {
		delete(cm.conns, conn)
		for conn, id := range cm.conns {
			if id > removeId {
				cm.conns[conn]--
			}
		}
		cm.currId--
	} else {
		fmt.Println("Did not find a connection to remove in the connection list.")
	}
}
