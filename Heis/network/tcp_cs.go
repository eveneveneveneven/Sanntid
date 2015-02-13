package network

import (
    "fmt"
    "log"
    "net"
    "encoding/gob"
)

type P struct {
	m, n int64
}

type Client struct {
	raddr 	string	
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

type Connection struct {
	conn    net.Conn
	encoder *gob.Encoder
	decoder *gob.Decoder
}

type Server struct {
	listener net.Listener
	conns []*Connection
}

func NewClient(raddr string, port string) *Client {
	var c Client
	c.raddr = raddr + ":" + port
	
	return &c
}

func (c *Client) ConnectServer(createGob bool) error {
	fmt.Println(c.raddr)
	conn, err := net.Dial("tcp", c.raddr)
	if err != nil {
		return err
	}
	c.conn = conn
	if createGob {
		c.encoder = gob.NewEncoder(conn)
		c.decoder = gob.NewDecoder(conn)
	}
	fmt.Println("Got connection!")
	p := &P{1, 2}
	c.encoder.Encode(p)
	c.conn.Close()
	return nil
}

func (c *Client) CleanUp() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func NewServer(port int) *Server {
	var s Server
	laddr := ":20011"
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = ln
	return &s
}

func (s *Server) handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &P{}
	dec.Decode(p)
	fmt.Printf("Received server : %+v\n", p)
	conn.Close()
}

func (s *Server) ListenServer() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Got faulty TCP connection server side")
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) CleanUp() {
	if s.listener != nil {
		s.listener.Close()
	}
}
