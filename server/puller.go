package server

func NewPuller(session *Session, pusher *Pusher) *Puller {
	return &Puller{
		pusher:  pusher,
		session: session,
	}
}

type Puller struct {
	session *Session
	pusher  *Pusher
}

func (c *Puller) Pause(pause bool) {

}