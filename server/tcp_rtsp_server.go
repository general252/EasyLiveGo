package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func NewTcpRtspServer(port int) *TcpRtspServer {
	return &TcpRtspServer{port: port}
}

type TcpRtspServer struct {
	listen *net.TCPListener
	wg     sync.WaitGroup
	port   int

	pushers sync.Map // path -> Pusher
}

func (c *TcpRtspServer) Start() error {
	var err error
	c.listen, err = net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: c.port,
		Zone: "",
	})
	if err != nil {
		return err
	}

	c.wg.Add(1)
	go c.loop()

	return nil
}

func (c *TcpRtspServer) Stop() {
	_ = c.listen.Close()

	c.wg.Wait()
}

func (c *TcpRtspServer) loop() {
	defer c.wg.Done()

	log.Printf("listen %v", c.listen.Addr())
	for {
		tcpConn, err := c.listen.AcceptTCP()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("new connection: %v", tcpConn.RemoteAddr())

		_ = tcpConn.SetReadBuffer(64 * 1024)
		_ = tcpConn.SetWriteBuffer(64 * 1024)

		session := NewSession(tcpConn, c)
		session.Start()
	}
}

func (c *TcpRtspServer) AddPusher(pusher *Pusher) {
	c.pushers.Store(pusher.Path, pusher)
}

func (c *TcpRtspServer) GetPusher(path string) (*Pusher, error) {
	obj, ok := c.pushers.Load(path)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	res, ok := obj.(*Pusher)
	if !ok || res == nil {
		return nil, fmt.Errorf("error")
	}

	return res, nil
}
