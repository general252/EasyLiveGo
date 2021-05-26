package server

import (
	"fmt"
	"log"
	"net"
)

func NewPuller(session *Session, pusher *Pusher) *Puller {
	c := &Puller{
		Id:      generalShortId(32),
		session: session,
		pusher:  pusher,
	}

	c.pusher.AddPuller(c)

	return c
}

type Puller struct {
	Id string

	session *Session
	pusher  *Pusher

	addrAudio     *net.UDPAddr
	addrAudioCtrl *net.UDPAddr
	addrVideo     *net.UDPAddr
	addrVideoCtrl *net.UDPAddr
}

func (c *Puller) Stop() {
	if c.pusher != nil {
		c.pusher.RemovePuller(c)
	}
}

func (c *Puller) Pause(pause bool) {

}

func (c *Puller) SetupAudio(host string, audioPort int, audioCtrlPort int) {
	var port int
	var err error

	if len(host) == 0 {
		log.Println("host is nil")
		return
	}

	port = audioPort
	c.addrAudio, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Println(err)
		return
	}

	port = audioCtrlPort
	c.addrAudioCtrl, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Puller) SetupVideo(host string, videoPort, videoCtrlPort int) {
	var port int
	var err error
	if len(host) == 0 {
		log.Println("host is nil")
		return
	}

	port = videoPort
	c.addrVideo, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Println(err)
		return
	}

	port = videoCtrlPort
	c.addrVideoCtrl, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%v:%v", host, port))
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Puller) handleRtp(rtp *RtpPack) {
	var addr *net.UDPAddr
	if rtp.Type == PacketTypeVideo {
		addr = c.addrVideo
	} else if rtp.Type == PacketTypeAudio {
		addr = c.addrAudio
	} else {
		log.Println("unknown packet type")
	}

	if addr != nil {
		n, err := rtp.Conn.WriteToUDP(rtp.Buffer.Bytes(), addr)
		if err != nil {
			log.Println(err)
		} else {
			_ = n
		}
	} else {
		log.Println("rtp addr is nil")
	}
}

func (c *Puller) handleRtcp(rtcp *RtcpPack) {
	var addr *net.UDPAddr
	if rtcp.Type == PacketTypeVideo {
		addr = c.addrVideoCtrl
	} else if rtcp.Type == PacketTypeAudio {
		addr = c.addrAudioCtrl
	} else {
		log.Println("unknown packet type")
	}

	if addr != nil {
		n, err := rtcp.Conn.WriteToUDP(rtcp.Buffer.Bytes(), addr)
		if err != nil {
			log.Println(err)
		} else {
			_ = n
		}
	} else {
		log.Println("rtcp addr is nil")
	}
}
