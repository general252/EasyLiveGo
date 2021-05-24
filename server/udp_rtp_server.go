package server

import (
	"log"
	"net"
	"sync"
)

func NewUdpRtpServer(port int) *UdpRtpServer {
	return &UdpRtpServer{port: port}
}

type UdpRtpServer struct {
	listen *net.UDPConn
	wg     sync.WaitGroup
	port   int
}

func (c *UdpRtpServer) Start() error {
	var err error
	c.listen, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: c.port,
		Zone: "",
	})
	if err != nil {
		return err
	}

	_ = c.listen.SetWriteBuffer(64 * 1024)
	_ = c.listen.SetReadBuffer(64 * 1024)

	c.wg.Add(1)
	go c.loop()

	return nil
}

func (c *UdpRtpServer) Stop() {
	_ = c.listen.Close()

	c.wg.Wait()
}

func (c *UdpRtpServer) loop() {
	defer c.wg.Done()

	buffer := make([]byte, 65535)
	var connList = make(map[string]string)
	for {
		n, addr, err := c.listen.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		msg := buffer[:n]

		_ = addr
		_ = msg

		if _, ok := connList[addr.String()]; !ok {
			connList[addr.String()] = addr.String()
			log.Printf("new rtp conn %v", addr)
		}
	}
}
