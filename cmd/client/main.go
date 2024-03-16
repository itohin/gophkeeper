package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/logger"
	"io"
	"os"
)

func main() {
	l := logger.NewLogger()

	shutdownCh := make(chan struct{})

	fingerPrint, err := makeFingerPrint()
	if err != nil {
		l.Fatal(err)
	}
	client, err := grpc.NewClient(fingerPrint, shutdownCh)
	if err != nil {
		l.Fatal(err)
	}
	defer client.Close()

	authUseCase := auth.NewAuth(client)

	p := prompt.NewPrompt()
	app := cli.NewCli(l, p, authUseCase, shutdownCh)

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
