package cli

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/router"
	"github.com/itohin/gophkeeper/pkg/logger"
)

type Auth interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	Verify(ctx context.Context, login, otp string) error
}

const (
	authMenu = "authMenu"
	register = "register"
	login    = "login"
	verify   = "verify"

	registerLabel = "Регистрация"
	loginLabel    = "Вход"
	verifyLabel   = "Подтвердить email"

	dataMenu = "dataMenu"
	getData  = "getData"
	addData  = "addData"

	addDataLabel = "Сохранить данные"
	getDataLabel = "Получить данные"
)

type Cli struct {
	router *router.Router
	log    logger.Logger
	prompt prompt.Prompter
	auth   Auth
}

func NewCli(
	logger logger.Logger,
	prompt prompt.Prompter,
	auth Auth,
) *Cli {
	cli := &Cli{
		log:    logger,
		prompt: prompt,
		auth:   auth,
	}

	cli.router = router.NewRouter(
		map[string]router.Command{
			authMenu: cli.authMenu,
			register: cli.register,
			login:    cli.login,
			verify:   cli.verify,
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
	c.log.Infof("requested action: %s", action)

	for {
		cmd, err := c.router.GetCommand(action)
		if err != nil {
			return err
		}
		action, err = cmd()
		if err != nil {
			return err
		}

		c.log.Infof("requested action: %s", action)
	}
}
