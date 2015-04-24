package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"

	"../types"
)

const (
	TCP_PORT = 30011
)

func readFromTCPConn(decoder *gob.Decoder, recieve chan *types.NetworkMessage, stop chan bool) {
	for {
		msg := &types.NetworkMessage{}
		if err := decoder.Decode(msg); err != nil {
			fmt.Printf("\t\t\x1b[31;1mError\x1b[0m |readFromTCPConn| [%v]\n", err)
			stop <- true
			return
		}
		recieve <- msg
	}
}

func sendToTCPConn(encoder *gob.Encoder, msg *types.NetworkMessage) error {
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

func createTCPHandler(conn *net.TCPConn, wakeRecieve, wakeSend chan *types.NetworkMessage,
	connEnd chan *net.TCPConn, terminate chan bool, wg *sync.WaitGroup) {

	fmt.Println("\t\tStarting new TCP handler!")
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	stop := make(chan bool, 1)
	recieve := make(chan *types.NetworkMessage)
	numErrorSend := 0
	go readFromTCPConn(decoder, recieve, stop)
	for {
		// prioritized channel to check
		select {
		case <-stop:
			connEnd <- conn
			conn.Close()
			return
		case <-terminate:
			conn.Close()
			return
		default:
		}

		select {
		case <-stop:
			connEnd <- conn
			conn.Close()
			return
		case <-terminate:
			conn.Close()
			return
		case recieveMsg := <-recieve:
			wakeRecieve <- recieveMsg
		case msgHolder := <-wakeSend:
			if err := sendToTCPConn(encoder, msgHolder); err != nil {
				fmt.Printf("\t\t\x1b[31;1mError\x1b[0m |createTCPHandler| [%v]\n", err)
				numErrorSend++
				if numErrorSend >= 5 {
					fmt.Println("\t\tFailed to send msg 5 times, stops connection")
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
	fmt.Println("\t\tStarting TCP listener!")
	laddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP:   net.ParseIP("localhost"),
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Printf("\t\t\x1b[31;1mError\x1b[0m |startTCPListener| [%v], quitting program\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("\t\t\x1b[31;1mError\x1b[0m |startTCPListener| [%v], continue listening\n", err)
			continue
		}
		newConn <- conn
	}
}

func createTCPConn(ip string) (*net.TCPConn, error) {
	fmt.Println("\t\tCreating TCP connection")
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
