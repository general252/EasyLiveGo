package server

import (
	"bytes"
	"net"
)

type SessionType int

const (
	SessionTypeUnknown SessionType = iota
	SessionTypePusher
	SessionTypePuller
)

type PacketType int

const (
	PacketTypeUnknown PacketType = iota
	PacketTypeVideo
	PacketTypeAudio
)

type RtpPack struct {
	Conn   *net.UDPConn
	Type   PacketType
	Buffer *bytes.Buffer
}

func NewRtpPack(pktType PacketType, data []byte, srcConn *net.UDPConn) *RtpPack {
	return &RtpPack{
		Conn:   srcConn,
		Type:   pktType,
		Buffer: bytes.NewBuffer(data),
	}
}

type RtcpPack struct {
	Conn   *net.UDPConn
	Type   PacketType
	Buffer *bytes.Buffer
}

func NewRtcpPack(pktType PacketType, data []byte, srcConn *net.UDPConn) *RtcpPack {
	return &RtcpPack{
		Conn:   srcConn,
		Type:   pktType,
		Buffer: bytes.NewBuffer(data),
	}
}
