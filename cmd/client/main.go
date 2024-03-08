package main

import (
	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	auth2 "github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/logger"
)

func main() {
	l := logger.NewLogger()
	p := prompt.NewPrompt()
	auth := auth2.NewAuth()

	app := cli.NewCli(l, p, auth)

	err := app.Run()
	if err != nil {
		l.Fatal(err)
	}
}
