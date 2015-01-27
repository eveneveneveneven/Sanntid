package main

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
)

type P struct {
    X, Y int64
}

func main() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    encoder := gob.NewEncoder(conn)
    p := &P{123, 321}
    encoder.Encode(p)
    conn.Close()
    fmt.Println("done");
}