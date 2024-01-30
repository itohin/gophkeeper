package cli

import (
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/code_generator"
	"github.com/itohin/gophkeeper/pkg/logger"
	"github.com/itohin/gophkeeper/pkg/mailer"
)

type Cli struct {
	log           logger.Logger
	prompt        prompt.Prompter
	mailer        mailer.Mailer
	codeGenerator code_generator.Generator
	auth          auth.Auth
}

func NewCli(
	logger logger.Logger,
	prompt prompt.Prompter,
	mailer mailer.Mailer,
	codeGen code_generator.Generator,
	auth auth.Auth,
) *Cli {
	return &Cli{
		log:           logger,
		prompt:        prompt,
		mailer:        mailer,
		codeGenerator: codeGen,
		auth:          auth,
	}
}
