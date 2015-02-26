package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

func readFromTCPConn(conn *net.TCPConn, recieve chan *networkMessage, stop chan bool) {
	decoder := gob.NewDecoder(conn)
	for {
		msg := &networkMessage{}
		if err := decoder.Decode(msg); err != nil {
			fmt.Printf("Some error %v\n", err)
			stop <- true
			return
		}
		recieve <- msg
	}
}

func sendToTCPConn(conn *net.TCPConn, msg *networkMessage) error {
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

func createTCPHandler(conn *net.TCPConn, wakeRecieve, wakeSend chan *networkMessage, connEnd chan *net.TCPConn, wg *sync.WaitGroup) {
	stop := make(chan bool)
	recieve := make(chan *networkMessage)
	go readFromTCPConn(conn, recieve, stop)
	numErrorSend := 0
	for {
		// prioritized channel to check
		select {
		case <-stop:
			connEnd <- conn
			conn.Close()
			return
		default:
		}

		select {
		case <-stop:
			connEnd <- conn
			conn.Close()
			return
		case recieveMsg := <-recieve:
			wakeRecieve <- recieveMsg
		case sendMsg := <-wakeSend:
			if err := sendToTCPConn(conn, sendMsg); err != nil {
				fmt.Printf("Some error %v\n", err)
				numErrorSend++
				if numErrorSend >= 5 {
					fmt.Println("Failed to send msg 5 times, stops connection")
					connEnd <- conn
					conn.Close()
					wg.Done()
					return
				}
			} else {
				numErrorSend = 0
			}
			wg.Done()
		}
	}
}

func startTCPListener(newConn chan *net.TCPConn) {
	laddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP:   net.ParseIP("localhost"),
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Printf("Some error %v, quitting program\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("Some error %v, continue listening\n", err)
			continue
		}
		newConn <- conn
	}
}

func createConnTCP(ip string) (*net.TCPConn, error) {
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
