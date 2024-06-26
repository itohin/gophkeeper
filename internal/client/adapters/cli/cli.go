package cli

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/router"
	"github.com/itohin/gophkeeper/internal/client/entities"
	errors2 "github.com/itohin/gophkeeper/pkg/errors"
)

type Auth interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
	Verify(ctx context.Context, login, otp string) error
	Logout(ctx context.Context) error
}

type Secrets interface {
	CreateSecret(ctx context.Context, secret *entities.Secret) error
	GetSecrets(ctx context.Context) (map[string]*entities.Secret, error)
	GetSecret(ctx context.Context, id string) (*entities.Secret, error)
	DeleteSecret(ctx context.Context, id string) error
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
	dataMenu         = "dataMenu"
	getData          = "getData"
	addData          = "addData"
	deleteData       = "deleteData"
	addText          = "addText"
	addPassword      = "addPassword"
	addBinary        = "addBinary"
	addCard          = "addCard"
	saveBinaryToDisk = "saveBinaryToDisk"
	showData         = "showData"

	addDataLabel          = "Сохранить данные"
	getDataLabel          = "Получить данные"
	deleteDataLabel       = "Удалить данные"
	addTextLabel          = "Текстовые данные"
	addBinaryLabel        = "Бинарные данные"
	addCardLabel          = "Данные банковской карты"
	saveBinaryToDiskLabel = "Сохранить на диске"
	addPasswordLabel      = "Данные для входа(логин/пароль)"

	comeBackLabel = "Вернуться назад"
)

type Cli struct {
	router     *router.Router
	prompt     prompt.Prompter
	auth       Auth
	secrets    Secrets
	shutdownCh chan struct{}
	errorCh    chan error
}

func NewCli(
	prompt prompt.Prompter,
	auth Auth,
	secrets Secrets,
	shutdownCh chan struct{},
	errorCh chan error,
) *Cli {
	cli := &Cli{
		prompt:     prompt,
		auth:       auth,
		secrets:    secrets,
		shutdownCh: shutdownCh,
		errorCh:    errorCh,
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
			addCard:     cli.addCard,
			addPassword: cli.addPassword,
			addBinary:   cli.addBinary,
			showData:    cli.showData,
			deleteData:  cli.deleteData,
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
		case err := <-c.errorCh:
			fmt.Println("\n\n", err.Error())
			action, err = c.Call(authMenu)
		default:
			action, err = c.Call(action)
			if err != nil {
				if errors.As(err, &domainError) {
					fmt.Println("\n\n", err.Error())
					action, err = c.Call(authMenu)
				} else {
					return err
				}
			}
		}
	}
}

func (c *Cli) Call(action string) (result string, err error) {
	if len(action) < 1 {
		return "", fmt.Errorf("empty action name")
	}
	actionData := strings.Split(action, "/")
	cmd, err := c.router.GetCommand(actionData[0])
	if err != nil {
		return "", err
	}
	if reflect.TypeOf(cmd).Kind() != reflect.Func {
		return "", fmt.Errorf("requested action %v not a func", cmd)
	}

	params := actionData[1:]
	f := reflect.ValueOf(cmd)
	if len(params) != f.Type().NumIn() {
		err = fmt.Errorf("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value
	res = f.Call(in)
	if len(res) < 2 {
		err = fmt.Errorf("not enough parameters returned from action %v", cmd)
		return
	}
	if !res[1].IsNil() {
		return "", res[1].Interface().(error)
	}
	result = res[0].String()
	return
}
