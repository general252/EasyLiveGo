package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func NewTcpRtSpServer(port int) *TcpRtSpServer {
	return &TcpRtSpServer{port: port}
}

type TcpRtSpServer struct {
	listen *net.TCPListener
	wg     sync.WaitGroup
	port   int

	sessionList sync.Map // sessionId -> Session
}

func (c *TcpRtSpServer) Start() error {
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

func (c *TcpRtSpServer) Stop() {
	_ = c.listen.Close()

	c.wg.Wait()
}

func (c *TcpRtSpServer) loop() {
	defer c.wg.Done()

	log.Printf("skt %v", c.listen.Addr())
	for {
		tcpConn, err := c.listen.AcceptTCP()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("new connection: %v", tcpConn.RemoteAddr())

		_ = tcpConn.SetReadBuffer(512 * 1024)
		_ = tcpConn.SetWriteBuffer(512 * 1024)

		session := NewSession(tcpConn, c)
		c.sessionList.Store(session.Id, session)

		session.Start()
	}
}

func (c *TcpRtSpServer) GetPusher(path string) (*Pusher, error) {
	var res *Pusher

	c.SessionRange(func(session *Session) bool {
		if session.Type == SessionTypePusher {
			if session.pusher.Path == path {
				res = session.pusher
				return false
			}
		}

		return true
	})

	if res == nil {
		return nil, fmt.Errorf("not found")
	}

	return res, nil
}

func (c *TcpRtSpServer) GetPusherByAddr(host string, port int) (*Pusher, error) {
	var res *Pusher

	c.SessionRange(func(session *Session) bool {
		if session.Type == SessionTypePusher {
			if session.Host() == host {
				if session.VPort == port || session.APort == port || session.VControlPort == port || session.AControlPort == port {
					res = session.pusher
					return false
				}
			}
		}

		return true
	})

	if res == nil {
		return nil, fmt.Errorf("not found")
	}

	return res, nil
}

func (c *TcpRtSpServer) SessionRange(f func(*Session) bool) {
	if c == nil || f == nil {
		return
	}

	c.sessionList.Range(func(key, value interface{}) bool {
		session, ok := value.(*Session)
		if !ok || session == nil {
			return true
		}

		return f(session)
	})
}
