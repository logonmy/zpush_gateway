package utils

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"zpush/conf"
)

var redisConn redis.Conn

func InitRedis() error {
	conn, err := redis.Dial("tcp", conf.Config().Redis.Address, redis.DialPassword(conf.Config().Redis.Password))
	if err != nil {
		log.Printf("dial redis error: %s\n", err.Error())
		return err
	}

	redisConn = conn
	return nil
}

func RedisConn() redis.Conn {
	return redisConn
}
