package server

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/general252/go-rtsp/rtsp"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func NewSession(conn *net.TCPConn, server *TcpRtSpServer) *Session {
	const networkBuffer = 200 * 1024
	connRW := bufio.NewReadWriter(bufio.NewReaderSize(conn, networkBuffer), bufio.NewWriterSize(conn, networkBuffer))
	return &Session{
		Id:        generalSessionId(),
		tcpServer: server,
		conn:      conn,
		connRW:    connRW,
		Type:      SessionTypeUnknown,
	}
}

type Session struct {
	Id string

	tcpServer *TcpRtSpServer

	conn   *net.TCPConn
	connRW *bufio.ReadWriter

	Type   SessionType
	Url    string
	Path   string
	sdp    string // SDP
	sdpMap map[string]*rtsp.SDPInfo

	AControl string // audio control 示例: streamid=0
	VControl string // video control 示例: streamid=1
	ACodec   string // audio codec 示例: aac
	VCodec   string // video codec 示例: h264

	APort        int // client audio port
	AControlPort int // client audio control port
	VPort        int // client video port
	VControlPort int // client video control port

	pusher *Pusher // 推流流(或拉流对应的发流)
	puller *Puller // 拉流
}

func (c *Session) Start() {
	go c.loop()
}

func (c *Session) Stop() {
	c.conn.Close()
	c.tcpServer.sessionList.Delete(c.Id)
}

func (c *Session) Host() string {
	addr, ok := c.conn.RemoteAddr().(*net.TCPAddr)
	if ok && addr != nil {
		return addr.IP.String()
	}

	return ""
}

func (c *Session) loop() {
	defer c.Stop()

	reqBuf := bytes.NewBuffer(nil)
	for {
		line, isPrefix, err := c.connRW.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		reqBuf.Write(line)
		if !isPrefix {
			reqBuf.WriteString("\r\n")
		}

		if len(line) == 0 {
			req := rtsp.NewRequest(reqBuf.String())
			if req == nil {
				log.Printf("parse RTSP request fail")
				break
			}

			contentLen := req.GetContentLength()
			if contentLen > 0 {
				bodyBuf := make([]byte, contentLen)
				if n, err := io.ReadFull(c.connRW, bodyBuf); err != nil {
					log.Println(err)
					return
				} else if n != contentLen {
					log.Printf("read rtsp request body failed, expect size[%d], got size[%d]", contentLen, n)
					return
				}
				req.Body = string(bodyBuf)
			}

			c.handleRequest(req)
			reqBuf.Reset()
		}
	}
}

func (c *Session) handleRequest(req *rtsp.Request) {

	log.Printf("收到数据:\n%v", req)

	res := rtsp.NewResponse(200, "OK", req.Header["CSeq"], c.Id, "")
	defer func() {
		if p := recover(); p != nil {
			log.Printf("handleRequest err ocurs:%v", p)
			res.StatusCode = 500
			res.Status = fmt.Sprintf("Inner Server Error, %v", p)
		}

		outBytes := []byte(res.String())
		if _, err := c.connRW.Write(outBytes); err != nil {
			log.Println(err)
			c.Stop()
		}
		if err := c.connRW.Flush(); err != nil {
			log.Println(err)
		}

		log.Printf("发送数据:\n%v", res)

		switch req.Method {
		case "PLAY", "RECORD":
		case "TEARDOWN":
			{
				c.Stop()
				return
			}
		}

		if res.StatusCode != 200 && res.StatusCode != 401 {
			log.Printf("Response request error[%d]. stop session.", res.StatusCode)
			c.Stop()
		}

	}()

	var announce = func() {
		c.Type = SessionTypePusher
		c.Url = req.URL

		reqURL, err := url.Parse(req.URL)
		if err != nil {
			res.StatusCode = 500
			res.Status = "Invalid URL"
			return
		}
		c.Path = reqURL.Path

		c.sdp = req.Body
		c.sdpMap = rtsp.ParseSDP(c.sdp)

		sdp, ok := c.sdpMap["audio"]
		if ok {
			c.AControl = sdp.Control
			c.ACodec = sdp.Codec
			log.Printf("audio codec[%s]\n", c.ACodec)
		}

		sdp, ok = c.sdpMap["video"]
		if ok {
			c.VControl = sdp.Control
			c.VCodec = sdp.Codec
			log.Printf("video codec[%s]\n", c.VCodec)
		}

		c.pusher = NewPusher(c, c.Path)
	}

	var setup = func() {
		ts := req.Header["Transport"]
		// control字段可能是`stream=1`字样，也可能是rtsp://...字样。即control可能是url的path，也可能是整个url
		// 例1：
		// a=control:streamid=1
		// 例2：
		// a=control:rtsp://192.168.1.64/trackID=1
		// 例3：
		// a=control:?ctype=video
		setupUrl, err := url.Parse(req.URL)
		if err != nil {
			res.StatusCode = 500
			res.Status = "Invalid URL"
			return
		}
		if setupUrl.Port() == "" {
			setupUrl.Host = fmt.Sprintf("%s:554", setupUrl.Host)
		}
		setupPath := setupUrl.String()

		mTcp := regexp.MustCompile("interleaved=(\\d+)(-(\\d+))?")
		mUdp := regexp.MustCompile("client_port=(\\d+)(-(\\d+))?")

		tcpMatches := mTcp.FindStringSubmatch(ts)
		if tcpMatches != nil {
			// tcp transport
			res.StatusCode = 500
			res.Status = "not support"
			return
		}

		udpMatches := mUdp.FindStringSubmatch(ts)
		if udpMatches == nil {
			res.StatusCode = 500
			res.Status = "Invalid URL"
			return
		}
		// udp transport

		vPath := c.VControl
		aPath := c.AControl

		if setupPath == aPath || aPath != "" && strings.LastIndex(setupPath, aPath) == len(setupPath)-len(aPath) {
			// audio
			c.APort, _ = strconv.Atoi(udpMatches[1])
			c.AControlPort, _ = strconv.Atoi(udpMatches[3])

			switch c.Type {
			case SessionTypePuller:
				//
				if c.puller != nil {
					c.puller.SetupAudio(c.Host(), c.APort, c.AControlPort)
				}
			case SessionTypePusher:
				//
				tss := strings.Split(ts, ";")
				idx := -1
				for i, val := range tss {
					if val == udpMatches[0] {
						idx = i
					}
				}
				tail := append([]string{}, tss[idx+1:]...)
				tss = append(tss[:idx+1], fmt.Sprintf("server_port=%d-%v", 5020, 5021))
				tss = append(tss, tail...)
				ts = strings.Join(tss, ";")
			}
		} else if setupPath == vPath || vPath != "" && strings.LastIndex(setupPath, vPath) == len(setupPath)-len(vPath) {
			// video
			c.VPort, _ = strconv.Atoi(udpMatches[1])
			c.VControlPort, _ = strconv.Atoi(udpMatches[3])

			switch c.Type {
			case SessionTypePuller:
				//
				if c.puller != nil {
					c.puller.SetupVideo(c.Host(), c.VPort, c.VControlPort)
				}
			case SessionTypePusher:
				//
				tss := strings.Split(ts, ";")
				idx := -1
				for i, val := range tss {
					if val == udpMatches[0] {
						idx = i
					}
				}
				tail := append([]string{}, tss[idx+1:]...)
				tss = append(tss[:idx+1], fmt.Sprintf("server_port=%d-%v", 5020, 5021))
				tss = append(tss, tail...)
				ts = strings.Join(tss, ";")
			}
		}
		res.Header["Transport"] = ts
	}

	var describe = func() {
		c.Type = SessionTypePuller

		c.Url = req.URL

		reqURL, err := url.Parse(req.URL)
		if err != nil {
			res.StatusCode = 500
			res.Status = "Invalid URL"
			return
		}
		c.Path = reqURL.Path

		pusher, err := c.tcpServer.GetPusher(c.Path)
		if err != nil {
			res.StatusCode = 404
			res.Status = "NOT FOUND"
			return
		}

		c.pusher = pusher
		c.AControl = pusher.AControl()
		c.VControl = pusher.VControl()
		c.ACodec = pusher.ACodec()
		c.VCodec = pusher.VCodec()

		c.puller = NewPuller(c, pusher)
		res.SetBody(pusher.session.sdp)
	}

	switch req.Method {
	case rtsp.OPTIONS:
		res.Header["Public"] = "DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE, OPTIONS, ANNOUNCE, RECORD"
	case rtsp.ANNOUNCE:
		announce()
	case rtsp.DESCRIBE:
		describe()
	case rtsp.SETUP:
		setup()
	case rtsp.PLAY:
		if c.pusher == nil {
			res.StatusCode = 500
			res.Status = "Error Status"
			return
		}
		res.Header["Range"] = req.Header["Range"]
	case rtsp.RECORD:
		if c.pusher == nil {
			res.StatusCode = 500
			res.Status = "Error Status"
			return
		}
	case rtsp.PAUSE:
		if c.puller == nil {
			res.StatusCode = 500
			res.Status = "Error Status"
			return
		}
		c.puller.Pause(true)
	}
}

// generalSessionId create session id
func generalSessionId() string {
	dest := make([]byte, 16)
	n, err := rand.Read(dest)
	if err != nil {
		return fmt.Sprintf("%v", rand.Uint64())
	}

	h := md5.New()
	h.Write(dest[:n])
	return hex.EncodeToString(h.Sum(nil))[8:24]
}
