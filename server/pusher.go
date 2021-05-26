package server

import (
	"sync"
)

func NewPusher(session *Session, path string) *Pusher {
	r := &Pusher{
		session:    session,
		Path:       path,
		pullerList: make(map[string]*Puller),
	}

	return r
}

type RtpHandle func(*RtpPack)
type RtcpHandle func(*RtcpPack)

type Pusher struct {
	session *Session
	Path    string

	pullerList map[string]*Puller
	pullerMux  sync.Mutex
}

func (c *Pusher) Stop() {
	var pullers []*Puller
	c.lock(func() {
		for _, puller := range c.pullerList {
			pullers = append(pullers, puller)
		}
	})

	c.pullerList = make(map[string]*Puller)

	for _, puller := range pullers {
		puller.Stop()
	}
}

func (c *Pusher) lock(f func()) {
	c.pullerMux.Lock()
	defer c.pullerMux.Unlock()

	f()
}

func (c *Pusher) Range(f func(*Puller) bool) {
	c.lock(func() {
		for _, puller := range c.pullerList {
			if f(puller) == false {
				return
			}
		}
	})
}

func (c *Pusher) AddPuller(pull *Puller) {
	c.lock(func() {
		c.pullerList[pull.Id] = pull
	})
}

func (c *Pusher) RemovePuller(pull *Puller) {
	c.lock(func() {
		delete(c.pullerList, pull.Id)
	})
}

func (c *Pusher) HandleRtp(rtp *RtpPack) {
	c.lock(func() {
		for _, puller := range c.pullerList {
			puller.handleRtp(rtp)
		}
	})
}

func (c *Pusher) HandleRtcp(rtcp *RtcpPack) {
	c.lock(func() {
		for _, puller := range c.pullerList {
			puller.handleRtcp(rtcp)
		}
	})
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
