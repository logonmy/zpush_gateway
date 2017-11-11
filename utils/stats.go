package utils

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"sync/atomic"
	"time"
	"zpush/conf"
)

type ServerStats struct {
	StartTime time.Time
	MsgIn     uint64
	MsgOut    uint64
}

var stats ServerStats

func (s *ServerStats) RecordStart() {
	s.StartTime = time.Now()

	go s.ReportStats()
}

func (s *ServerStats) RecordMsgIn() {
	atomic.AddUint64(&s.MsgIn, 1)
}

func (s *ServerStats) RecordMsgOut() {
	atomic.AddUint64(&s.MsgOut, 1)
}

func Stats() *ServerStats {
	return &stats
}

func (s *ServerStats) ReportStats(){
	tick := time.NewTicker(time.Second * 5)
	for{
		select{
		case <- tick.C:
			go s.Report()
		}
	}
}

func (s *ServerStats) Report() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     conf.Config().Influxdb.Address,
		Username: conf.Config().Influxdb.Username,
		Password: conf.Config().Influxdb.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conf.Config().Influxdb.DB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}


	fields := map[string]interface{}{
		"msg_in":   float32(s.MsgIn),
		"msg_out": float32(s.MsgOut),
	}

	pt, err := client.NewPoint("zpush_msg", nil, fields)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	log.Println("report stats to influxdb success")
}
