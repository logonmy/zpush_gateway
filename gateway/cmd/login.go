package cmd

import (
	"log"
	msg "zpush/gateway/message"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"fmt"
	"github.com/pkg/errors"
)


func onLogin(cmdCode uint16, packet []byte) (proto.Message, error) {
	var req msg.LoginReq
	err := proto.Unmarshal(packet, &req)
	if err != nil{
		log.Printf("unmarshal client msg error: %s\n", err.Error())
		return nil, errors.New("unmarshal client msg error")
	}

	userid := 10001 + rand.Intn(50000)

	resp := &msg.LoginResp{
		Userid: int32(userid),
		Token: fmt.Sprintf("faketoken%d", userid),
	}
	return resp, nil
}
