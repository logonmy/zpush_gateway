package main

import (
	"flag"
	"log"

	"context"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"google.golang.org/grpc"
	"time"
	"zpush/conf"
	"zpush/gateway"
	"zpush/http"
	"zpush/rpc"
	pb "zpush/rpc/protocol"
	"zpush/utils"
)

var (
	configFile = flag.String("config", "./config.toml", "config file path")
)

func test_client() {
	time.Sleep(time.Second * 5)

	userid := 10086

	key := fmt.Sprintf("user:gw:%d", userid)

	gateway_addr, err := redis.String(utils.RedisConn().Do("GET", key))
	if err != nil {
		log.Fatalln("redis error")
	}

	conn, err := grpc.Dial(gateway_addr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewMsgServiceClient(conn)

	for {
		_, err := client.Send(context.Background(), &pb.MsgReq{Userid: 10086, Content: "this is msg from logic"})
		if err != nil {
			log.Printf("RPC CALL error: %s\n", err.Error())
			return
		}

		time.Sleep(time.Second * 1)
	}

}

func main() {
	log.Println("start zpush gateway")

	flag.Parse()
	err := conf.Parse(*configFile)
	if err != nil {
		log.Fatalf("parse config file error: %s\n", err.Error())
	}

	utils.Stats().RecordStart()
	
	go http.StartHTTPServer()
	go rpc.StartRPCServer()

	//go test_client()

	gateway.StartTCPServer()

	log.Println("zpush exit.")
}
