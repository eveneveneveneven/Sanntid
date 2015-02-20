package main

import (
	"fmt"
	"net"
	"time"
	"strconv"
	"os/exec"
)

func main() {
	count := 0
	p := make([]byte, 1024)
	searching := true
	
	laddr, _ := net.ResolveUDPAddr("udp", "localhost:20011")
	ln, _ := net.ListenUDP("udp", laddr)

	for searching {
		timeout := make(chan bool)
		connection := make(chan bool)
		go func() {
			time.Sleep(200 * time.Millisecond)
			timeout <- true
		}()
		go func() {
			n, _, err := ln.ReadFromUDP(p)
			if err == nil {
				connection <- true
				val, _ := strconv.Atoi(string(p[:n]))
				count = val
			}
		}()
		select {
		case <-timeout:
			searching = false
			count += 1
			ln.Close()
		case <-connection:
		}
	}
	
	cmd := exec.Command("gnome-terminal", "-e", "./counter")
	cmd.Output()
	
	raddr, _ := net.ResolveUDPAddr("udp", ":20011")
	conn, _ := net.DialUDP("udp", nil, raddr)
	defer conn.Close()
	
	for {
		fmt.Printf("Count: %v\n", count)
		fmt.Fprintf(conn, strconv.Itoa(count))
		count += 1
		time.Sleep(100 * time.Millisecond)
	}
}
