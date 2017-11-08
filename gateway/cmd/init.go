package cmd

import (
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"log"
)

const (
	HEADER_LEN = 2 + 4
)

type PacketHeader struct {
	Cmd     uint16
	BodyLen uint32
}

func (this *PacketHeader) Unmarshal(b []byte) error {
	n := 0

	this.Cmd = binary.BigEndian.Uint16(b[n:])
	n += 2

	this.BodyLen = binary.BigEndian.Uint32(b[n:])
	n += 4

	return nil
}

func ParsePacketHeader(packetHeaderBuf []byte) (*PacketHeader, error) {
	var packetHeader PacketHeader

	err := packetHeader.Unmarshal(packetHeaderBuf)
	if err != nil {
		log.Printf("unmarshal packet header error: %s\n", err.Error())
		return nil, errors.New("unmarshal packet header error")
	}

	return &packetHeader, nil
}

type CMD_HANDLER func(uint16, []byte) (proto.Message, error)

var (
	cmdHandler = map[uint16]CMD_HANDLER{
		1: onLogin,
	}
)

func DispatchCmd(packetHeader *PacketHeader, packetBodyBuf []byte) (proto.Message, error) {
	handler, ok := cmdHandler[packetHeader.Cmd]
	if !ok {
		log.Printf("unsupport cmd code: %d\n", packetHeader.Cmd)
		return nil, nil
	}

	respMsg, err := handler(packetHeader.Cmd, packetBodyBuf)
	if err != nil {
		return nil, err
	}

	return respMsg, nil
}
