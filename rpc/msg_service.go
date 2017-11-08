package rpc

import (
	"context"
	"zpush/gateway"
	protocol "zpush/rpc/protocol"
)

type MsgService struct {
}

func NewMsgService() *MsgService {
	return &MsgService{}
}

func (s *MsgService) Send(ctx context.Context, req *protocol.MsgReq) (*protocol.MsgResp, error) {
	gateway.SendMsg(int(req.Userid), []byte(req.Content))
	return &protocol.MsgResp{}, nil
}
