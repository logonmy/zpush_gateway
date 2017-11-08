package gateway

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"sync"
	"time"
	"zpush/conf"
	"zpush/utils"
)

var (
	GWServer *Gateway
)

func SendMsg(userid int, msg []byte) {
	if GWServer == nil {
		log.Println("GWServer is nil, can't send msg")
		return
	}

	GWServer.sendMsg(userid, msg)
}

type Gateway struct {
	sessions           map[int]*Session
	userSessions       map[int]*Session
	sessionsLocker     sync.RWMutex
	userSessionsLocker sync.RWMutex
	sessionId          int
}

func NewGateway() *Gateway {
	return &Gateway{
		sessions:     make(map[int]*Session),
		userSessions: make(map[int]*Session),
		sessionId:    1,
	}
}

func (this *Gateway) sendMsg(userid int, msg []byte) error {
	this.sessionsLocker.RLock()
	userSession, ok := this.userSessions[userid]
	this.sessionsLocker.RUnlock()

	if !ok {
		log.Printf("未找到userid: %d的session对象\n", userid)
		return nil
	}

	userSession.SendMsg(msg)
	return nil
}

func (this *Gateway) gatewayTimer() {
	timer := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-timer.C:
			this.userSessionsLocker.RLock()
			//log.Println("dump userSessions map")
			for userid, session := range this.userSessions {
				log.Printf("userid: %d —> session: %v\n", userid, session)
			}
			this.userSessionsLocker.RUnlock()
		}
	}
}

func (this *Gateway) start() {
	log.Println("gatway server started")

	err := utils.InitRedis()
	if err != nil {
		log.Fatalf("init redis error: %s\n", err.Error())
	}
	log.Println("init redis sucessful")

	go this.gatewayTimer()

	listener, err := net.Listen("tcp", conf.Config().Server.TcpBind)
	utils.PanicIfError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept connection error: %s\n", err.Error())
			continue
		}

		session := NewSession(this, conn, this.sessionId)
		go session.Process()

		this.updateSessions(session)
	}
}

func (this *Gateway) updateUserSession(session *Session) {
	this.userSessionsLocker.Lock()
	defer this.userSessionsLocker.Unlock()

	this.userSessions[session.userId] = session

	// 将用户的session信息写入db
	key := fmt.Sprintf("user:gw:%d", session.userId)

	_, err := redis.String(utils.RedisConn().Do("SET", key, conf.Config().Server.RpcBind))
	if err != nil {
		log.Printf("set user session to redis error: %s\n", err.Error())
	}
}

func (this *Gateway) updateSessions(session *Session) {
	this.sessionsLocker.Lock()
	defer this.sessionsLocker.Unlock()

	this.sessions[this.sessionId] = session
	this.sessionId += 1
}

func (this *Gateway) DeleteSession(session *Session) {
	this.userSessionsLocker.Lock()
	defer this.userSessionsLocker.Unlock()

	delete(this.userSessions, session.userId)

	key := fmt.Sprintf("user:gw:%d", session.userId)
	_, err := utils.RedisConn().Do("DEL", key)
	if err != nil {
		log.Printf("delete user session from redis error: %s\n", err.Error())
	}

}

func (this *Gateway) registerSession(sessionId int, session *Session) {

}

func StartTCPServer() {
	GWServer = NewGateway()
	GWServer.start()
}
