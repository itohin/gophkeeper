package cli

import (
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/router"
	"github.com/itohin/gophkeeper/internal/client/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/code_generator"
	"github.com/itohin/gophkeeper/pkg/logger"
	"github.com/itohin/gophkeeper/pkg/mailer"
)

const (
	authMenu = "authMenu"
	register = "register"
	login    = "login"
	dataMenu = "dataMenu"
	getData  = "getData"
	addData  = "addData"

	registerLabel = "Регистрация"
	loginLabel    = "Вход"
	addDataLabel  = "Сохранить данные"
	getDataLabel  = "Получить данные"
)

type Cli struct {
	router        *router.Router
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
	cli := &Cli{
		log:           logger,
		prompt:        prompt,
		mailer:        mailer,
		codeGenerator: codeGen,
		auth:          auth,
	}

	cli.router = router.NewRouter(
		map[string]router.Command{
			authMenu: cli.authMenu,
			register: cli.register,
			login:    cli.login,
			dataMenu: cli.dataMenu,
			getData:  cli.getData,
			addData:  cli.addData,
		},
	)

	return cli
}

func (c *Cli) Run() error {
	action, err := c.authMenu()
	if err != nil {
		return err
	}

	for {
		cmd, err := c.router.GetCommand(action)
		if err != nil {
			return err
		}
		action, err = cmd()
		if err != nil {
			return err
		}
	}
}
