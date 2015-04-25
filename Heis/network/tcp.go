package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"../types"
)

const (
	TCP_PORT = 30011

	READ_DEADLINE  = 500 // milliseconds
	WRITE_DEADLINE = 500 // milliseconds
)

func readFromTCPConn(conn *net.TCPConn, decoder *gob.Decoder,
	recieve chan *types.NetworkMessage, stop chan bool) {
	for {
		msg := &types.NetworkMessage{}
		conn.SetReadDeadline(time.Now().Add(READ_DEADLINE * time.Millisecond))
		if err := decoder.Decode(msg); err != nil {
			fmt.Printf("\t\t\x1b[31;1mError\x1b[0m |readFromTCPConn| [%v]\n", err)
			stop <- true
			return
		}
		recieve <- msg
	}
}

func sendToTCPConn(conn *net.TCPConn, encoder *gob.Encoder,
	msg *types.NetworkMessage) error {
	conn.SetWriteDeadline(time.Now().Add(WRITE_DEADLINE * time.Millisecond))
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

func createTCPHandler(conn *net.TCPConn, wakeRecieve, wakeSend chan *types.NetworkMessage,
	connEnd chan *net.TCPConn, terminate chan bool, wg *sync.WaitGroup) {

	fmt.Println("\x1b[34;1m::: Start New TCP Handler :::\x1b[0m")

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	stop := make(chan bool, 1)
	recieve := make(chan *types.NetworkMessage)
	numErrorSend := 0
	go readFromTCPConn(conn, decoder, recieve, stop)
	for {
		// prioritized channel to check
		select {
		case <-stop:
			conn.Close()
			connEnd <- conn
			return
		case <-terminate:
			conn.Close()
			return
		default:
		}

		select {
		case <-stop:
			conn.Close()
			connEnd <- conn
			return
		case <-terminate:
			conn.Close()
			return
		case recieveMsg := <-recieve:
			wakeRecieve <- recieveMsg
		case msgHolder := <-wakeSend:
			if err := sendToTCPConn(conn, encoder, msgHolder); err != nil {
				fmt.Printf("\t\x1b[31;1mError\x1b[0m |createTCPHandler| [%v]\n", err)
				numErrorSend++
				if numErrorSend >= 5 {
					fmt.Println("\t\x1b[31;1mError\x1b[0m |createTCPHandler| [tFailed to send msg 5 times, stops connection]")
					connEnd <- conn
					conn.Close()
					wg.Done()
					return
				}
			} else {
				numErrorSend = 0
			}
			wg.Done()
			wg.Wait()
		}
	}
}

func startTCPListener(newConn chan *net.TCPConn) {
	fmt.Println("\x1b[34;1m::: Start TCP Listener :::\x1b[0m")

	laddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP:   net.ParseIP("localhost"),
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |startTCPListener| [%v], quitting program\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("\t\x1b[31;1mError\x1b[0m |startTCPListener| [%v], continue listening\n", err)
			continue
		}
		newConn <- conn
	}
}

func createTCPConn(ip string) (*net.TCPConn, error) {
	fmt.Println("\x1b[36;1m::: Creating New TCP Connection :::\x1b[0m")
	raddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
