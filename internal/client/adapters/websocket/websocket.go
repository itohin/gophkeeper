package websocket

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/itohin/gophkeeper/internal/server/events"
	"io"
	"log"
	"sync"
)

type WSListener struct {
	url         string
	fingerPrint string
	shutdownCh  chan struct{}
}

func NewWSListener(url, fingerPrint string, shutdownCh chan struct{}) *WSListener {
	return &WSListener{
		url:         url,
		fingerPrint: fingerPrint,
		shutdownCh:  shutdownCh,
	}
}

func (w *WSListener) Listen(ctx context.Context, userID string) error {
	dialer := ws.Dialer{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, _, err := dialer.Dial(ctx, fmt.Sprintf("%s?finger_print=%s&user_id=%s", w.url, w.fingerPrint, userID))
	if err != nil {
		return fmt.Errorf("failed to dial the websocket server: %v", err)
	}
	defer conn.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go w.receiveMsg(ctx, wg, conn)

	<-w.shutdownCh

	if err := closeConnection(conn); err != nil {
		return fmt.Errorf("failed to close ws connection: %v", err)
	}
	wg.Wait()

	log.Println("connection closed")

	return nil
}

func (w *WSListener) receiveMsg(ctx context.Context, wg *sync.WaitGroup, conn io.ReadWriteCloser) {
	defer wg.Done()

	for {
		select {
		case <-w.shutdownCh:
			return
		default:
		}

		h, r, err := wsutil.NextReader(conn, ws.StateClientSide)
		log.Println("ws read msg", h, r, err)
		if err != nil {
			log.Printf("failed to read a frame: %v", err)
			continue
		}
		switch h.OpCode {
		case ws.OpClose:
			return
		case ws.OpText:
			msg, err := io.ReadAll(r)
			if err != nil {
				log.Printf("failed to read a message: %v", err)
				continue
			}
			var s events.SecretEvent
			err = json.Unmarshal(msg, &s)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			log.Printf("received message: %v", s)
		default:
			log.Printf("unexpected frame with opcode %d", h.OpCode)
		}
	}
}

func closeConnection(conn io.ReadWriteCloser) error {
	if err := wsutil.WriteClientMessage(conn, ws.OpClose, nil); err != nil {
		return fmt.Errorf("failed to write a close frame: %v", err)
	}
	return nil
}
