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
	skt  *net.UDPConn
	wg   sync.WaitGroup
	port int
}

func (c *UdpRtpServer) Start() error {
	var err error
	c.skt, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: c.port,
		Zone: "",
	})
	if err != nil {
		return err
	}

	_ = c.skt.SetWriteBuffer(512 * 1024)
	_ = c.skt.SetReadBuffer(512 * 1024)

	c.wg.Add(1)
	go c.loop()

	return nil
}

func (c *UdpRtpServer) Stop() {
	_ = c.skt.Close()

	c.wg.Wait()
}

func (c *UdpRtpServer) loop() {
	defer c.wg.Done()

	var connList = make(map[string]*Pusher)

	buffer := make([]byte, 65535)
	for {
		n, addr, err := c.skt.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		msg := buffer[:n]

		_ = addr
		_ = msg

		pusher, ok := connList[addr.String()]
		if !ok {
			pusher, _ := DefaultApp.GetTcpServer().GetPusherByAddr(addr.IP.String(), addr.Port)
			connList[addr.String()] = pusher
			log.Printf("new rtp conn %v, %v", addr, pusher.GetPath())
		} else {
			var pktType = PacketTypeUnknown
			if pusher.session.APort == addr.Port {
				pktType = PacketTypeAudio
			} else if pusher.session.VPort == addr.Port {
				pktType = PacketTypeVideo
			}

			var pkt = NewRtpPack(pktType, msg, c.skt)
			pusher.HandleRtp(pkt)
		}
	}
}
