package cmd

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"log"
	msg "zpush/gateway/message"
)

func onHeartbeat(packet []byte) (proto.Message, error) {
	var req msg.HBReq
	err := proto.Unmarshal(packet, &req)
	if err != nil {
		log.Printf("unmarshal client msg error: %s\n", err.Error())
		return nil, errors.New("unmarshal client msg error")
	}

	return &msg.HBResp{}, nil
}
