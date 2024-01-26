package main

import (
	"github.com/itohin/gophkeeper/internal/adapters/cli"
	"github.com/itohin/gophkeeper/internal/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/infrastructure/code_generator"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/itohin/gophkeeper/internal/infrastructure/mailer"
)

func main() {
	l := logger.NewLogger()
	p := prompt.NewPrompt()
	m := mailer.NewSMTPMailer("from@gmail.com", "", "localhost", "1025")
	c := code_generator.NewCodeGenerator()

	app := cli.NewCli(l, p, m, c)

	err := app.Auth()
	if err != nil {
		l.Fatal(err)
	}
}
