package websocket

import (
	"io"
	"sync"
)

type Client struct {
	mx       *sync.Mutex
	conn     io.ReadWriteCloser
	id       string
	deviceID string
	done     chan struct{}
}

func NewClient(id, deviceID string, conn io.ReadWriteCloser) *Client {
	return &Client{
		mx:       &sync.Mutex{},
		conn:     conn,
		id:       id,
		deviceID: deviceID,
		done:     make(chan struct{}),
	}
}
