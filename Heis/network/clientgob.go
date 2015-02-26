package network

import (
    "fmt"
    "net"
    "encoding/gob"
    "time"
)

type P struct {
    M, N string
}

func client() {
	fmt.Println("start client");
	for {
		time.Sleep(1 * time.Second)
		raddr, _ := net.ResolveTCPAddr("tcp", "129.241.187.147:20011")
		conn, err := net.DialTCP("tcp", nil, raddr)
		if err != nil {
		    fmt.Println("got no connection")
		    continue
		}
		encoder := gob.NewEncoder(conn)
		p := &P{"yo", "man"}
		encoder.Encode(p)
		conn.Close()
	}
    fmt.Println("done");
}

func handleConnection(conn *net.TCPConn) {
    dec := gob.NewDecoder(conn)
    p := &P{}
    dec.Decode(p)
    fmt.Printf("Received : %+v\n", p);
}

func server() {
	fmt.Println("start server");
	laddr, _ := net.ResolveTCPAddr("tcp", ":20011")
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.AcceptTCP() // this blocks until connection or error
		if err != nil {
		    fmt.Println("Something went wrong")
		    continue
		}
		// a goroutine handles conn so that the loop can accept other connections
		go handleConnection(conn)
	}
}


func main() {
	ch := make(chan int)
    go client()
    go server()
    <-ch
}
