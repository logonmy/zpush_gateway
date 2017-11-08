package utils

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
	"zpush/conf"
	"fmt"
)

const (
	ROOT_PATH = "/gateway"
)

var zkConn *zk.Conn

func InitZK() error {
	conn, _, err := zk.Connect([]string{conf.Config().ZK.Address}, time.Second*5)
	if err != nil {
		log.Printf("connect to zookeeper error: %s\n", err.Error())
		return err
	}

	zkConn = conn
	return nil
}

func ensureRootPath() error {
	_, err := zkConn.Create(ROOT_PATH, []byte{1}, 0, zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		log.Println("create znode err", err, "path=", ROOT_PATH)
		return err
	}

	return nil
}

func RegisterGateway(nodeId string, addr string) error {
	err := ensureRootPath()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", ROOT_PATH, nodeId)

	_, err = zkConn.CreateProtectedEphemeralSequential(path, []byte(addr), zk.WorldACL(zk.PermAll))
	if err != nil && err != zk.ErrNodeExists {
		log.Println("create server znode err, path=", addr, err)
		return err
	}

	return nil
}
