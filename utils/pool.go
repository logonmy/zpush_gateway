package utils

import (
	"bufio"
	"net"
	"sync"
)

var (
	pool sync.Pool
)

func ReleaseBufReader(reader *bufio.Reader) {
	if reader == nil {
		return
	}

	pool.Put(reader)
}

func GetBufReader(conn net.Conn) (reader *bufio.Reader) {
	poolItem := pool.Get()
	if poolItem == nil {
		reader = bufio.NewReader(conn)
	} else {
		reader = poolItem.(*bufio.Reader)
		reader.Reset(conn)
	}

	return
}
