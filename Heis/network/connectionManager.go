package network

import (
	"fmt"
	"net"
	"sync"
)

type connManager struct {
	masterIP string
	currId   int
	conns    map[*net.TCPConn]int

	wakeRecieve chan *networkMessage
	wakeSend    chan *networkMessage
	newConn     chan *net.TCPConn
	connEnd     chan *net.TCPConn

	hubRecieve chan *networkMessage
	hubSend    chan *networkMessage

	wg *sync.WaitGroup
}

func NewConnManager(hbRec, hbSend chan *networkMessage) *connManager {
	var cm connManager

	cm.masterIP = ""
	cm.currId = 1
	cm.conns = make(map[*net.TCPConn]int)

	cm.wakeRecieve = make(chan *networkMessage, 20) // buffer for messages recieved
	cm.wakeSend = make(chan *networkMessage)
	cm.newConn = make(chan *net.TCPConn)
	cm.connEnd = make(chan *net.TCPConn)

	cm.hubRecieve = hbRec
	cm.hubSend = hbSend

	return &cm
}

func (cm *connManager) run() {
	go startTCPListener(cm.newConn)
	for {
		// prioritized channels to check
		select {
		case conn := <-cm.connEnd:
			cm.removeConnection(conn)
			continue
		case conn := <-cm.newConn:
			cm.addConnection(conn)
			go createTCPHandler(conn, cm.wakeRecieve, cm.wakeSend, cm.connEnd, cm.wg)
			continue
		default:
		}

		select {
		case conn := <-cm.connEnd:
			cm.removeConnection(conn)
		case conn := <-cm.newConn:
			cm.addConnection(conn)
			go createTCPHandler(conn, cm.wakeRecieve, cm.wakeSend, cm.connEnd, cm.wg)
		case recieveMsg := <-cm.wakeRecieve:
			cm.hubRecieve <- recieveMsg
		case sendMsg := <-cm.hubSend:
			numConns := len(cm.conns)
			if numConns > 0 {
				cm.wg.Add(numConns)
				for i := 0; i < numConns; i++ {
					cm.wakeRecieve <- sendMsg
				}
				cm.wg.Wait()
			}
		}
	}
}

func (cm *connManager) connectToNetwork(masterIP string) error {
	cm.masterIP = masterIP
	conn, err := createConnTCP(cm.masterIP)
	if err != nil {
		return err
	}
	cm.addConnection(conn)
	go createTCPHandler(conn, cm.wakeRecieve, cm.wakeSend, cm.connEnd, cm.wg)
	return nil
}

func (cm *connManager) addConnection(conn *net.TCPConn) {
	cm.conns[conn] = cm.currId
	cm.currId++
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
