package main

import (
	"github.com/itohin/gophkeeper/internal/client/adapters/cli"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	auth2 "github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/code_generator"
	"github.com/itohin/gophkeeper/pkg/logger"
	"github.com/itohin/gophkeeper/pkg/mailer"
)

func main() {
	l := logger.NewLogger()
	p := prompt.NewPrompt()
	m := mailer.NewSMTPMailer("from@gmail.com", "", "localhost", "1025")
	c := code_generator.NewCodeGenerator()
	auth := auth2.NewAuth()

	app := cli.NewCli(l, p, m, c, auth)

	err := app.Run()
	if err != nil {
		l.Fatal(err)
	}
}
