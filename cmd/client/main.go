package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/internal/client/usecases/secrets"
	"github.com/itohin/gophkeeper/pkg/jwt"
	"github.com/itohin/gophkeeper/pkg/logger"
	"io"
	"os"
	"time"
)

func main() {
	l := logger.NewLogger()

	shutdownCh := make(chan struct{})

	fingerPrint, err := makeFingerPrint()
	if err != nil {
		l.Fatal(err)
	}
	//TODO: maybe change to jwtClaimManager
	jwtGen, err := jwt.NewJWTGOManager("secret", 60*time.Second, 360*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	token := entities.NewToken(jwtGen)
	client, err := grpc.NewClient(fingerPrint, token, shutdownCh)
	if err != nil {
		l.Fatal(err)
	}
	defer client.Close()

	authUseCase := auth.NewAuth(client)
	secretsUseCase := secrets.NewSecrets(client)

	p := prompt.NewPrompt()
	app := cli.NewCli(l, p, authUseCase, secretsUseCase, shutdownCh)

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
