package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/client/adapters/storage"
	"github.com/itohin/gophkeeper/internal/client/adapters/websocket"
	conf "github.com/itohin/gophkeeper/internal/client/config"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/internal/client/usecases/secrets"
	"github.com/itohin/gophkeeper/pkg/jwt"
)

func main() {
	cfg := conf.ReadConfig()

	fmt.Println("cfg: ", cfg.JWT.Signature)

	shutdownCh := make(chan struct{})
	authCh := make(chan string, 1)
	errorCh := make(chan error)

	fingerPrint, err := makeFingerPrint()
	if err != nil {
		log.Fatal(err)
	}

	jwtGen, err := jwt.NewJWTGOManager(cfg.JWT.Signature, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	if err != nil {
		log.Fatal(err)
	}
	token := entities.NewToken(jwtGen)
	hydrator := storage.NewSecretsHydrator()
	client, err := grpc.NewClient(fingerPrint, token, shutdownCh, hydrator, cfg.GRPC.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	memoryStorage := storage.NewMemoryStorage()
	authUseCase := auth.NewAuth(client, authCh)
	secretsUseCase := secrets.NewSecrets(client, memoryStorage)

	ws := websocket.NewWSListener(
		fmt.Sprintf("wss://%s/connect", cfg.WebSocket.ServerAddress),
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
				ctx, cancel := context.WithTimeout(context.Background(), cfg.WebSocket.ConnectionTimeout)
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
	app := cli.NewCli(p, authUseCase, secretsUseCase, shutdownCh, errorCh)

	err = app.Start()
	if err != nil {
		log.Fatal(err)
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
