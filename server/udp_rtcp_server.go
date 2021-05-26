package server

import (
	"log"
	"net"
	"sync"
)

func NewUdpRtcpServer(port int) *UdpRtcpServer {
	return &UdpRtcpServer{port: port}
}

type UdpRtcpServer struct {
	skt  *net.UDPConn
	wg   sync.WaitGroup
	port int
}

func (c *UdpRtcpServer) Start() error {
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

func (c *UdpRtcpServer) Stop() {
	_ = c.skt.Close()

	c.wg.Wait()
}

func (c *UdpRtcpServer) loop() {
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
			log.Printf("new rtcp conn %v, %v", addr, pusher.GetPath())
		} else {
			var pktType = PacketTypeUnknown
			if pusher.session.AControlPort == addr.Port {
				pktType = PacketTypeAudio
			} else if pusher.session.VControlPort == addr.Port {
				pktType = PacketTypeVideo
			}

			var pkt = NewRtcpPack(pktType, msg, c.skt)
			pusher.HandleRtcp(pkt)
		}
		// log.Printf("rtcp %v %v", addr, string(msg))
	}
}
