package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
)

func readFromTCPConn(decoder *gob.Decoder, recieve chan *networkMessage, stop chan bool) {
	for {
		msg := &networkMessage{}
		if err := decoder.Decode(msg); err != nil {
			fmt.Printf("\t\tError |readFromTCPConn| [%v]\n", err)
			stop <- true
			return
		}
		recieve <- msg
	}
}

func sendToTCPConn(encoder *gob.Encoder, msg *networkMessage) error {
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	return nil
}

func createTCPHandler(conn *net.TCPConn, wakeRecieve, wakeSend chan *networkMessage,
	connEnd chan *net.TCPConn, wg *sync.WaitGroup) {

	fmt.Println("\t\tStarting new TCP handler!")
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)
	stop := make(chan bool)
	recieve := make(chan *networkMessage)
	numErrorSend := 0
	go readFromTCPConn(decoder, recieve, stop)
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
		case msgHolder := <-wakeSend:
			if err := sendToTCPConn(encoder, msgHolder); err != nil {
				fmt.Printf("\t\tError |createTCPHandler| [%v]\n", err)
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
		fmt.Printf("\t\tError |startTCPListener| [%v], quitting program\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Printf("\t\tError |startTCPListener| [%v], continue listening\n", err)
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
