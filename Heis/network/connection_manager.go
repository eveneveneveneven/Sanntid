package network

import (
	"fmt"
	"net"
	"sync"

	"../types"
)

const (
	BUFFER_MSG_RECIEVED = 20
)

type connection struct {
	id        int
	sendMsg   chan *types.NetworkMessage
	terminate chan bool
}

type connManager struct {
	masterIP string
	currId   int
	conns    map[*net.TCPConn]*connection

	wakeRecieve chan *types.NetworkMessage
	newConn     chan *net.TCPConn
	connEnd     chan *net.TCPConn

	hubRecieve chan *types.NetworkMessage
	hubSend    chan *types.NetworkMessage

	wg *sync.WaitGroup
}

func newConnManager(hbRec, hbSend chan *types.NetworkMessage) *connManager {
	return &connManager{
		masterIP: "",
		currId:   1,
		conns:    make(map[*net.TCPConn]*connection),

		// buffer for messages recieved
		wakeRecieve: make(chan *types.NetworkMessage, BUFFER_MSG_RECIEVED),
		newConn:     make(chan *net.TCPConn),
		connEnd:     make(chan *net.TCPConn),

		hubRecieve: hbRec,
		hubSend:    hbSend,

		wg: new(sync.WaitGroup),
	}
}

func (cm *connManager) run() {
	fmt.Println("\tStarting connection manager!")
	go startTCPListener(cm.newConn)
	for {
		// prioritized channels to check
		select {
		case conn := <-cm.connEnd:
			cm.removeConnection(conn)
			continue
		case conn := <-cm.newConn:
			cm.addConnection(conn)
			continue
		default:
		}

		select {
		case conn := <-cm.connEnd:
			cm.removeConnection(conn)
		case conn := <-cm.newConn:
			cm.addConnection(conn)
		case recieveMsg := <-cm.wakeRecieve:
			cm.hubRecieve <- recieveMsg
		case sendMsg := <-cm.hubSend:
			numConns := len(cm.conns)
			if numConns > 0 {
				fmt.Println("starting sending")
				cm.wg.Add(numConns)
				for _, c := range cm.conns {
					msgHolder := new(types.NetworkMessage)
					types.DeepCopy(msgHolder, sendMsg)
					msgHolder.Id = c.id
					c.sendMsg <- msgHolder
				}
				cm.wg.Wait()
				fmt.Println("sending done")
			}
		}
	}
}

func (cm *connManager) connectToNetwork(masterIP string) error {
	fmt.Printf("\tConnecting to network, Master ip:%v\n", masterIP)
	cm.masterIP = masterIP
	conn, err := createTCPConn(cm.masterIP)
	if err != nil {
		return err
	}
	cm.addConnection(conn)
	return nil
}

func (cm *connManager) addConnection(conn *net.TCPConn) {
	fmt.Printf("\tAdding connection [%v] with id %v\n", conn, cm.currId)
	c := &connection{
		id:        cm.currId,
		sendMsg:   make(chan *types.NetworkMessage),
		terminate: make(chan bool, 1),
	}
	cm.conns[conn] = c
	cm.currId++
	go createTCPHandler(conn, cm.wakeRecieve, c.sendMsg, cm.connEnd, c.terminate, cm.wg)
}

func (cm *connManager) removeConnection(conn *net.TCPConn) {
	if removeConn, ok := cm.conns[conn]; ok {
		fmt.Printf("\tRemoving connection [%v] with id %v\n", conn, removeConn.id)
		delete(cm.conns, conn)
		for conn, c := range cm.conns {
			if c.id > removeConn.id {
				cm.conns[conn].id--
			}
		}
		cm.currId--
	} else {
		fmt.Println("\t\t\x1b[31;1mError\x1b[0m |cm.removeConnection|",
			"[Did not find a connection to remove in the connection list]")
	}
}

func (cm *connManager) resetConnections() {
	for conn, c := range cm.conns {
		c.terminate <- true
		cm.removeConnection(conn)
	}
}
