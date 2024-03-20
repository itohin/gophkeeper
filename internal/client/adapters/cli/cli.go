package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/router"
	"github.com/itohin/gophkeeper/internal/client/entities"
	errors2 "github.com/itohin/gophkeeper/pkg/errors"
	"github.com/itohin/gophkeeper/pkg/logger"
)

type Auth interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	Verify(ctx context.Context, login, otp string) error
	Logout(ctx context.Context) error
}

type Secrets interface {
	CreateText(ctx context.Context, secret *entities.Secret, text string) error
	CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error
}

const (
	//роутинг
	//auth
	authMenu = "authMenu"
	register = "register"
	login    = "login"
	verify   = "verify"
	logout   = "logout"

	registerLabel = "Регистрация"
	loginLabel    = "Вход"
	verifyLabel   = "Подтвердить email"
	logoutLabel   = "Завершить работу"

	//data
	dataMenu    = "dataMenu"
	getData     = "getData"
	addData     = "addData"
	addText     = "addText"
	addPassword = "addPassword"

	addDataLabel     = "Сохранить данные"
	getDataLabel     = "Получить данные"
	addTextLabel     = "Текстовые данные"
	addPasswordLabel = "Данные для входа(логин/пароль)"

	comeBackLabel = "Вернуться назад"
)

type Cli struct {
	router     *router.Router
	log        logger.Logger
	prompt     prompt.Prompter
	auth       Auth
	secrets    Secrets
	shutdownCh chan struct{}
}

func NewCli(
	logger logger.Logger,
	prompt prompt.Prompter,
	auth Auth,
	secrets Secrets,
	shutdownCh chan struct{},
) *Cli {
	cli := &Cli{
		log:        logger,
		prompt:     prompt,
		auth:       auth,
		secrets:    secrets,
		shutdownCh: shutdownCh,
	}

	cli.router = router.NewRouter(
		map[string]router.Command{
			authMenu:    cli.authMenu,
			register:    cli.register,
			login:       cli.login,
			verify:      cli.verify,
			logout:      cli.logout,
			dataMenu:    cli.dataMenu,
			getData:     cli.getData,
			addData:     cli.addData,
			addText:     cli.addText,
			addPassword: cli.addPassword,
		},
	)

	return cli
}

func (c *Cli) Start() error {
	var domainError *errors2.DomainError
	action, err := c.authMenu()
	if err != nil {
		return err
	}

	for {
		select {
		case <-c.shutdownCh:
			return nil
		default:
			cmd, err := c.router.GetCommand(action)
			if err != nil {
				return err
			}
			action, err = cmd()
			if err != nil {
				if errors.As(err, &domainError) {
					fmt.Println("\n\n", err.Error())
					action = authMenu
				} else {
					return err
				}
			}
		}
	}
}
