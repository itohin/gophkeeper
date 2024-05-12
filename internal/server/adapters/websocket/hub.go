package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/itohin/gophkeeper/pkg/events"
)

type devicesMap map[string]*Client
type clientsMap map[string]devicesMap

type Hub struct {
	mx             *sync.RWMutex
	clients        clientsMap
	secretEventsCh chan *events.SecretEvent
}

func NewHub(secretEventsCh chan *events.SecretEvent) *Hub {
	return &Hub{
		mx:             &sync.RWMutex{},
		clients:        make(clientsMap),
		secretEventsCh: secretEventsCh,
	}
}

func (h *Hub) Connect(conn io.ReadWriteCloser, clientID, deviceID string) (<-chan struct{}, error) {
	c := h.createClient(conn, clientID, deviceID)
	err := h.addClient(clientID, deviceID, c)
	if err != nil {
		return nil, err
	}

	go h.handleBroadcast()
	go h.handleCloseClient(c)
	return c.done, nil
}

func (h *Hub) createClient(conn io.ReadWriteCloser, clientID, deviceID string) *Client {
	return NewClient(clientID, deviceID, conn)
}

func (h *Hub) handleBroadcast() {
	log.Println("start broadcast")
	for {
		select {
		case s := <-h.secretEventsCh:
			log.Printf("message received: %v", s)
			message, err := json.Marshal(s)
			if err != nil {
				log.Printf("failed to marshal data: %v", err)
			}

			err = h.broadcast(message)
			if err != nil {
				log.Printf("failed to broadcast message: %v", err)
			}
		default:

		}
	}

}

func (h *Hub) handleCloseClient(c *Client) {
	defer func(h *Hub, c *Client) {
		err := h.removeClient(c)
		if err != nil {
			log.Println(err)
		}
	}(h, c)
	for {
		hd, _, err := wsutil.NextReader(c.conn, ws.StateServerSide)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("failed to read a frame from the user id %s deviceId %s: %v", c.id, c.deviceID, err)
			continue
		}
		switch hd.OpCode {
		case ws.OpClose:
			if err := wsutil.WriteServerMessage(c.conn, ws.OpClose, nil); err != nil {
				log.Printf("failed to write a close frame: %v", err)
			}
			return
		default:
			log.Printf("unexpected frame with opcode %d", hd.OpCode)
		}
	}
}

func (h *Hub) broadcast(message []byte) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	for _, devices := range h.clients {
		for _, c := range devices {
			if err := wsutil.WriteServerText(c.conn, message); err != nil {
				return fmt.Errorf("failed to send a message to the client id %s deviceId %s: %v", c.id, c.deviceID, err)
			}
		}
	}
	return nil
}

func (h *Hub) addClient(clientID, deviceID string, c *Client) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	devices, ok := h.clients[clientID]
	if !ok {
		devices = make(devicesMap)
	}
	if _, ok := devices[deviceID]; ok {
		return fmt.Errorf("client id %s, deviceID %s already connected", clientID, deviceID)
	}
	devices[deviceID] = c
	h.clients[clientID] = devices
	return nil
}

func (h *Hub) removeClient(c *Client) error {
	h.mx.Lock()
	defer h.mx.Unlock()

	devices, ok := h.clients[c.id]
	if !ok {
		close(c.done)
		return fmt.Errorf("client id %s, deviceID %s not found", c.id, c.deviceID)
	}

	delete(devices, c.deviceID)
	if len(h.clients[c.id]) == 0 {
		delete(h.clients, c.id)
	}
	close(c.done)

	return nil
}
