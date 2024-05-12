package websocket

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/events"
)

type SecretsHolder interface {
	SaveSecret(ctx context.Context, secret *entities.Secret) error
	DeleteSecret(ctx context.Context, id string) error
}

type SecretHydrator interface {
	FromSecretEvent(event *events.SecretEvent) (*entities.Secret, error)
}

type WSListener struct {
	url           string
	fingerPrint   string
	shutdownCh    chan struct{}
	errorCh       chan error
	secretsHolder SecretsHolder
	hydrator      SecretHydrator
}

func NewWSListener(
	url, fingerPrint string,
	shutdownCh chan struct{},
	errorCh chan error,
	secretsHolder SecretsHolder,
	hydrator SecretHydrator,
) *WSListener {
	return &WSListener{
		url:           url,
		fingerPrint:   fingerPrint,
		shutdownCh:    shutdownCh,
		errorCh:       errorCh,
		secretsHolder: secretsHolder,
		hydrator:      hydrator,
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
			var event events.SecretEvent
			err = json.Unmarshal(msg, &event)
			if err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			err = w.handleEvent(ctx, event)
			if err != nil {
				log.Printf("failed to handle secret: %v", err)
				w.errorCh <- err
			}
		default:
			log.Printf("unexpected frame with opcode %d", h.OpCode)
		}
	}
}

func (w *WSListener) handleEvent(ctx context.Context, event events.SecretEvent) error {
	if event.EventType == events.TypeDeleted {
		err := w.secretsHolder.DeleteSecret(ctx, event.Secret.ID)
		if err != nil {
			return err
		}

		return nil
	}

	secret, err := w.hydrator.FromSecretEvent(&event)
	if err != nil {
		return err
	}

	return w.secretsHolder.SaveSecret(ctx, secret)
}

func closeConnection(conn io.ReadWriteCloser) error {
	if err := wsutil.WriteClientMessage(conn, ws.OpClose, nil); err != nil {
		return fmt.Errorf("failed to write a close frame: %v", err)
	}
	return nil
}
