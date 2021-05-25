package server

import (
	"fmt"
	"log"
	"net"
)

func NewPuller(session *Session, pusher *Pusher) *Puller {
	r := &Puller{
		pusher:  pusher,
		session: session,
	}

	pusher.AddRtpHandle(r.handleRtp)
	pusher.AddRtcpHandle(r.handleRtcp)

	return r
}

type Puller struct {
	session *Session
	pusher  *Pusher

	socketAudio *net.UDPConn
	socketVideo *net.UDPConn
}

func (c *Puller) Pause(pause bool) {

}

func (c *Puller) SetupAudio() {
	port := c.session.APort

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", c.session.Host, port))
	if err != nil {
		log.Println(err)
		return
	}

	c.socketAudio, err = net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Puller) SetupVideo() {
	port := c.session.VPort

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", c.session.Host, port))
	if err != nil {
		log.Println(err)
		return
	}

	c.socketVideo, err = net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Puller) handleRtp(rtp *RtpPack) {
	if rtp.Type == PacketTypeVideo {
		_, _ = c.socketVideo.Write(rtp.Buffer.Bytes())
	} else if rtp.Type == PacketTypeAudio {
		_, _ = c.socketAudio.Write(rtp.Buffer.Bytes())
	} else {
		log.Println("unknown packet type")
	}
}

func (c *Puller) handleRtcp(rtp *RtcpPack) {
	if rtp.Type == PacketTypeVideo {
		_, _ = c.socketVideo.Write(rtp.Buffer.Bytes())
	} else if rtp.Type == PacketTypeAudio {
		_, _ = c.socketAudio.Write(rtp.Buffer.Bytes())
	} else {
		log.Println("unknown packet type")
	}
}
