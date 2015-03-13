package network

import (
	"fmt"
	"net"
	"sync"

	"../types"
)

type connManager struct {
	masterIP string
	currId   int
	conns    map[*net.TCPConn]int

	wakeRecieve chan *types.NetworkMessage
	wakeSend    chan *types.NetworkMessage
	newConn     chan *net.TCPConn
	connEnd     chan *net.TCPConn

	hubRecieve chan *types.NetworkMessage
	hubSend    chan *types.NetworkMessage

	wg *sync.WaitGroup
}

func NewConnManager(hbRec, hbSend chan *types.NetworkMessage) *connManager {
	return &connManager{
		masterIP: "",
		currId:   1,
		conns:    make(map[*net.TCPConn]int),

		// buffer for messages recieved
		wakeRecieve: make(chan *types.NetworkMessage, types.BUFFER_MSG_RECIEVED),
		wakeSend:    make(chan *types.NetworkMessage),
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
				for i := 1; i <= numConns; i++ {
					msgHolder := new(types.NetworkMessage)
					*msgHolder = *sendMsg
					msgHolder.Id = i
					cm.wakeSend <- msgHolder
				}
				cm.wg.Wait()
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
	go createTCPHandler(conn, cm.wakeRecieve, cm.wakeSend, cm.connEnd, cm.wg)
	return nil
}

func (cm *connManager) addConnection(conn *net.TCPConn) {
	fmt.Printf("\tAdding connection [%v] with id %v\n", conn, cm.currId)
	cm.conns[conn] = cm.currId
	cm.currId++
}

func (cm *connManager) removeConnection(conn *net.TCPConn) {
	if removeId, ok := cm.conns[conn]; ok {
		fmt.Printf("\tRemoving connection [%v] with id %v\n", conn, removeId)
		delete(cm.conns, conn)
		for conn, id := range cm.conns {
			if id > removeId {
				cm.conns[conn]--
			}
		}
		cm.currId--
	} else {
		fmt.Println("\t\t\x1b[31;1mError\x1b[0m |cm.removeConnection|",
			"[Did not find a connection to remove in the connection list]")
	}
}
