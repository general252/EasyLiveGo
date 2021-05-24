package server

func NewPusher(session *Session, path string) *Pusher {
	return &Pusher{
		Path:    path,
		session: session,
	}
}

type Pusher struct {
	session *Session
	Path    string
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
