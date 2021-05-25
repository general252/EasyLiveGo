package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	DefaultApp = NewApp()
)

func NewApp() *App {
	return &App{}
}

type App struct {
	tcpServer     *TcpRtSpServer
	udpRtpServer  *UdpRtpServer
	udpRtcpServer *UdpRtcpServer
}

func (c *App) GetTcpServer() *TcpRtSpServer {
	return c.tcpServer
}

func (c *App) GetUdpServer() *UdpRtpServer {
	return c.udpRtpServer
}

func (c *App) Run() {

	rand.Seed(time.Now().Unix())

	// rtsp server
	c.tcpServer = NewTcpRtSpServer(554)
	if err := c.tcpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.tcpServer.Stop()

	// rtp server
	c.udpRtpServer = NewUdpRtpServer(5020)
	if err := c.udpRtpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.udpRtpServer.Stop()

	// rtcp server
	c.udpRtcpServer = NewUdpRtcpServer(5021)
	if err := c.udpRtcpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.udpRtcpServer.Stop()

	// http
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var sessionList []*Session
		c.GetTcpServer().SessionRange(func(session *Session) bool {
			sessionList = append(sessionList, session)
			return true
		})

		data, _ := json.MarshalIndent(sessionList, "", "  ")
		_, _ = w.Write(data)
	})

	// http server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
