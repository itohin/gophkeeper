package main

import (
	"github.com/itohin/gophkeeper/internal/adapters/cli"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/itohin/gophkeeper/internal/infrastructure/prompt"
)

func main() {
	l := logger.NewLogger()
	p := prompt.NewPrompt(l)
	app := cli.NewCli(l, p)

	err := app.Auth()
	if err != nil {
		l.Fatal(err)
	}
}
