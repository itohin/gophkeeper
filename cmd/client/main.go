package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/client/adapters/storage"
	"github.com/itohin/gophkeeper/internal/client/adapters/websocket"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/internal/client/usecases/secrets"
	"github.com/itohin/gophkeeper/pkg/jwt"
	"github.com/itohin/gophkeeper/pkg/logger"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	l := logger.NewLogger()

	shutdownCh := make(chan struct{})
	authCh := make(chan string, 1)
	errorCh := make(chan error)

	fingerPrint, err := makeFingerPrint()
	if err != nil {
		l.Fatal(err)
	}

	jwtGen, err := jwt.NewJWTGOManager("secret", 60*time.Second, 360*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	token := entities.NewToken(jwtGen)
	hydrator := storage.NewSecretsHydrator()
	client, err := grpc.NewClient(fingerPrint, token, shutdownCh, hydrator)
	if err != nil {
		l.Fatal(err)
	}
	defer client.Close()

	memoryStorage := storage.NewMemoryStorage()
	authUseCase := auth.NewAuth(client, authCh)
	secretsUseCase := secrets.NewSecrets(client, memoryStorage)

	wsPort := "7777"
	ws := websocket.NewWSListener(
		fmt.Sprintf("wss://:%s/connect", wsPort),
		fingerPrint,
		shutdownCh,
		errorCh,
		memoryStorage,
		hydrator,
	)

	go func() {
		for {
			select {
			case userID := <-authCh:
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*80)
				defer cancel()
				go func() {
					err := ws.Listen(ctx, userID)
					if err != nil {
						errorCh <- fmt.Errorf("ws listen error: %s", err)
					}
				}()
				err = secretsUseCase.SyncSecrets(context.Background())
				if err != nil {
					errorCh <- fmt.Errorf("не удалось синхронизировать данные: %v", err)
					log.Printf("ws listen error: %v", err)
				}
			}
		}
	}()

	p := prompt.NewPrompt()
	app := cli.NewCli(l, p, authUseCase, secretsUseCase, shutdownCh, errorCh)

	err = app.Start()
	if err != nil {
		l.Fatal(err)
	}
}

func makeFingerPrint() (string, error) {
	var fingerPrint string
	hostName, err := os.Hostname()
	if err != nil {
		return fingerPrint, err
	}
	hash := md5.New()
	_, err = io.WriteString(hash, hostName)
	if err != nil {
		return fingerPrint, err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
