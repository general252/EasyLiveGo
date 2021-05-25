package server

import "bytes"

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
	Type   PacketType
	Buffer *bytes.Buffer
}

func NewRtpPack(pktType PacketType, data []byte) *RtpPack {
	return &RtpPack{
		Type:   pktType,
		Buffer: bytes.NewBuffer(data),
	}
}

type RtcpPack struct {
	Type   PacketType
	Buffer *bytes.Buffer
}

func NewRtcpPack(pktType PacketType, data []byte) *RtcpPack {
	return &RtcpPack{
		Type:   pktType,
		Buffer: bytes.NewBuffer(data),
	}
}

