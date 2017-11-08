package rpc

import (
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"zpush/conf"
	protocol "zpush/rpc/protocol"
)

func StartRPCServer() {
	log.Println("rpc server started")

	lis, err := net.Listen("tcp", conf.Config().Server.RpcBind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	lis = netutil.LimitListener(lis, 10)

	s := grpc.NewServer()
	protocol.RegisterMsgServiceServer(s, NewMsgService())

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
