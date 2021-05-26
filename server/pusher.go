package server

import (
	"sync"
)

func NewPusher(session *Session, path string) *Pusher {
	r := &Pusher{
		Path:    path,
		session: session,
	}

	return r
}

type RtpHandle func(*RtpPack)
type RtcpHandle func(*RtcpPack)

type Pusher struct {
	session *Session
	Path    string

	rtpHandles    []RtpHandle
	rtpHandlesMux sync.Mutex

	rtcpHandles    []RtcpHandle
	rtcpHandlesMux sync.Mutex
}

func (c *Pusher) AddRtpHandle(rtpHandle RtpHandle) {
	c.rtpHandlesMux.Lock()
	defer c.rtpHandlesMux.Unlock()

	c.rtpHandles = append(c.rtpHandles, rtpHandle)
}

func (c *Pusher) AddRtcpHandle(rtcpHandle RtcpHandle) {
	c.rtcpHandlesMux.Lock()
	defer c.rtcpHandlesMux.Unlock()

	c.rtcpHandles = append(c.rtcpHandles, rtcpHandle)
}

func (c *Pusher) GetPath() string {
	if c == nil {
		return ""
	}
	return c.Path
}

func (c *Pusher) VCodec() string {
	if c.session != nil {
		return c.session.VCodec
	}
	return ""
}

func (c *Pusher) ACodec() string {
	if c.session != nil {
		return c.session.ACodec
	}
	return ""
}

func (c *Pusher) AControl() string {
	if c.session != nil {
		return c.session.AControl
	}
	return ""
}

func (c *Pusher) VControl() string {
	if c.session != nil {
		return c.session.VControl
	}
	return ""
}

func (c *Pusher) HandleRtp(rtp *RtpPack) {
	c.rtpHandlesMux.Lock()
	c.rtpHandlesMux.Unlock()

	for _, handle := range c.rtpHandles {
		handle(rtp)
	}
}

func (c *Pusher) HandleRtcp(rtcp *RtcpPack) {
	c.rtcpHandlesMux.Lock()
	c.rtcpHandlesMux.Unlock()

	for _, handle := range c.rtcpHandles {
		handle(rtcp)
	}
}
