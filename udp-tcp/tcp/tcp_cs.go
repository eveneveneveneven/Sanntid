package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func client(ch chan int) {
	fmt.Println("start client");
	
    conn, err := net.Dial("tcp", "129.241.187.136:33546")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    defer conn.Close()
   	
   	p := make([]byte, 1024)
	conn.Write([]byte("Connect to: 129.241.187.159:20011\x00")) // NB! remember 0-terminate!
	n, _ := conn.Read(p)
	fmt.Printf("We got back msg :: %s\n", p[:n])
    
    fmt.Println("client done");
    ch <- 0
}

func server(ch chan int) {
	fmt.Println("start server")
	
    ln, err := net.Listen("tcp", ":20011")
    if err != nil {
        log.Fatal(err)
    }
    defer ln.Close()
    
    conn, err := ln.Accept() // this blocks until connection or error
    if err != nil {
        log.Fatal("Listen error", err)
    }
    defer conn.Close()
    
    p := make([]byte, 1024)
    for {
    	n, _ := conn.Read(p)
    	fmt.Println("Got connection!")
		fmt.Printf("Received : %s\n", p[:n]);
		conn.Write([]byte("Send me more!\x00")) // NB! remember 0-terminate!
		time.Sleep(1 * time.Second)
    }
    
    ch <- 0
}


func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    
    go server(ch1)
    go client(ch2)
    
    <-ch1
    <-ch2
}
