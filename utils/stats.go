package utils

import (
	"time"
	"sync/atomic"
)

type ServerStats struct{
	StartTime time.Time
	MsgIn uint64
	MsgOut uint64
}


var stats ServerStats


func (s *ServerStats) RecordStart(){
	s.StartTime = time.Now()
}

func (s *ServerStats)RecordMsgIn(){
	atomic.AddUint64(&s.MsgIn, 1)
}

func (s *ServerStats) RecordMsgOut(){
	atomic.AddUint64(&s.MsgOut, 1)
}

func Stats() *ServerStats{
	return &stats
}