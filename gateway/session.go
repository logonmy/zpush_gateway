package gateway

import (
	"bufio"
	"encoding/binary"
	"github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	"zpush/gateway/cmd"
	msg "zpush/gateway/message"
	"zpush/utils"
)

type Session struct {
	conn        net.Conn
	userId      int
	sessionID   int
	reader      *bufio.Reader
	outPacketCh chan []byte
	connCloseCh chan bool
	gateway     *Gateway
}

func NewSession(gateway *Gateway, conn net.Conn, sessionId int) *Session {
	return &Session{
		gateway:     gateway,
		conn:        conn,
		sessionID:   sessionId,
		outPacketCh: make(chan []byte, 1000),
		connCloseCh: make(chan bool),
		reader:      utils.GetBufReader(conn),
	}
}

func (s *Session) SendMsg(msg []byte) {
	s.outPacketCh <- msg
}

func (s *Session) handleOutgoingPacket() {
	for {
		select {
		case <-s.connCloseCh:
			{
				return
			}
		case msg, ok := <-s.outPacketCh:
			if !ok {
				log.Println("outPacket channel closed")
				return
			}

			_, err := s.conn.Write(msg)
			if err != nil {
				log.Printf("send response to client error: %s\n", err.Error())
				return
			}

			utils.Stats().RecordMsgOut()
		}
	}
}

func (s *Session) removeSession() {
	s.gateway.DeleteSession(s)
}

func (s *Session) Process() {
	defer func() {
		if err := recover(); err != nil{
			log.Println(err)
		}
		s.connCloseCh <- true
		utils.ReleaseBufReader(s.reader)
		s.conn.Close()
	}()

	log.Printf("client: %s connected\n", s.conn.RemoteAddr().String())

	go s.handleOutgoingPacket()

	for {
		packetHeaderBuf := make([]byte, cmd.HEADER_LEN)
		_, err := io.ReadFull(s.reader, packetHeaderBuf)
		if err != nil {
			s.removeSession()
			return
		}

		packetHeader, err := cmd.ParsePacketHeader(packetHeaderBuf)
		if err != nil {
			log.Printf("unpack client header error: %s\n", err.Error())
			return
		}

		packetBodyBuf := make([]byte, packetHeader.BodyLen)

		_, err = io.ReadFull(s.reader, packetBodyBuf)
		if err != nil {
			log.Printf("io.ReadFull error: %s\n", err.Error())
			return
		}

		go s.handlerIncomingPacket(packetHeader, packetBodyBuf)
	}
}

func (s *Session) handlerIncomingPacket(packetHeader *cmd.PacketHeader, packetBodyBuf []byte) {
	utils.Stats().RecordMsgIn()

	respMsg, err := cmd.DispatchCmd(packetHeader, packetBodyBuf)
	if err != nil {
		// 表示系统处理失败
		s.outPacketCh <- utils.ErrSystem
		return
	}

	if packetHeader.Cmd == 1 {
		loginResp, _ := respMsg.(*msg.LoginResp)
		s.userId = int(loginResp.Userid)
		s.gateway.updateUserSession(s)
	}

	respBytes, err := proto.Marshal(respMsg)
	if err != nil {
		log.Println("marshal client msg error")
		return
	}

	buf := make([]byte, 2+4+len(respBytes))
	n := 0
	binary.BigEndian.PutUint16(buf[n:], 1)
	n += 2
	binary.BigEndian.PutUint32(buf[n:], uint32(len(respBytes)))
	n += 4

	copy(buf[n:], respBytes)
	s.outPacketCh <- buf
}
