package cli

import "C"
import (
	"errors"
	"fmt"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/itohin/gophkeeper/internal/infrastructure/prompt"
	"github.com/itohin/gophkeeper/internal/infrastructure/validator"
)

const (
	register = "Регистрация"
	login    = "Вход"
	addData  = "Сохранить данные"
	getData  = "Получить данные"
)

type Cli struct {
	log    logger.Logger
	prompt *prompt.Prompt
}

func NewCli(logger logger.Logger, prompt *prompt.Prompt) *Cli {
	return &Cli{
		log:    logger,
		prompt: prompt,
	}
}

func (c *Cli) Auth() error {
	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выполните вход или зарегистрируйтесь: "
	menuPrompt.ErrorMsg = "auth error"

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{login, register})
	if err != nil {
		return err
	}
	switch action {
	case login:
		return c.Login()
	case register:
		return c.Register()
	default:
		return fmt.Errorf("unknown action %s", action)
	}
}

func (c *Cli) Login() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин: "
	loginPrompt.ErrorMsg = "login error"
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(loginPrompt.ErrorMsg)
		}
		return nil
	}
	login, err := c.prompt.PromptGetInput(loginPrompt, validate)
	if err != nil {
		return err
	}
	passwordPrompt := prompt.PromptContent{
		Label:    "Введите пароль: ",
		ErrorMsg: "password error",
		Mask:     42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validate)
	if err != nil {
		return err
	}

	c.log.Info(login, password)
	//check credentials

	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выберите действие: "
	menuPrompt.ErrorMsg = "action error"

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{addData, getData})
	if err != nil {
		return err
	}
	c.log.Info(action)

	return nil
}

func (c *Cli) Register() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин(действующий email): "
	loginPrompt.ErrorMsg = "login error"
	login, err := c.prompt.PromptGetInput(loginPrompt, validator.ValidateEmail())
	if err != nil {
		return err
	}
	c.log.Info(login)
	return nil
}
