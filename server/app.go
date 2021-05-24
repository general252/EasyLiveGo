package server

import (
	"log"
	"net/http"
)

func NewApp() *App {
	return &App{}
}

type App struct {
	tcpServer     *TcpRtspServer
	udpRtpServer  *UdpRtpServer
	udpRtcpServer *UdpRtcpServer
}

func (c *App) GetTcpServer() *TcpRtspServer {
	return c.tcpServer
}

func (c *App) GetUdpServer() *UdpRtpServer {
	return c.udpRtpServer
}

func (c *App) Run() {
	c.tcpServer = NewTcpRtspServer(554)
	if err := c.tcpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.tcpServer.Stop()

	c.udpRtpServer = NewUdpRtpServer(5020)
	if err := c.udpRtpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.udpRtpServer.Stop()

	c.udpRtcpServer = NewUdpRtcpServer(5021)
	if err := c.udpRtcpServer.Start(); err != nil {
		log.Println(err)
		return
	}
	defer c.udpRtcpServer.Stop()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
