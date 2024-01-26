package cli

import "C"
import (
	"errors"
	"fmt"
	"github.com/itohin/gophkeeper/internal/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/infrastructure/code_generator"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/itohin/gophkeeper/internal/infrastructure/mailer"
	"github.com/itohin/gophkeeper/internal/infrastructure/validator"
)

const (
	register = "Регистрация"
	login    = "Вход"
	addData  = "Сохранить данные"
	getData  = "Получить данные"
)

type Cli struct {
	log           logger.Logger
	prompt        prompt.Prompter
	mailer        mailer.Mailer
	codeGenerator code_generator.Generator
}

func NewCli(logger logger.Logger, prompt prompt.Prompter, mailer mailer.Mailer, codeGen code_generator.Generator) *Cli {
	return &Cli{
		log:           logger,
		prompt:        prompt,
		mailer:        mailer,
		codeGenerator: codeGen,
	}
}

func (c *Cli) Auth() error {
	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выполните вход или зарегистрируйтесь: "

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{login, register})
	if err != nil {
		return err
	}
	switch action {
	case login:
		return c.login()
	case register:
		return c.register()
	default:
		return fmt.Errorf("unknown action %s", action)
	}
}

func (c *Cli) login() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин: "
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("error")
		}
		return nil
	}
	login, err := c.prompt.PromptGetInput(loginPrompt, validate)
	if err != nil {
		return err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validate)
	if err != nil {
		return err
	}

	c.log.Info(login, password)
	//check credentials

	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выберите действие: "

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{addData, getData})
	if err != nil {
		return err
	}
	c.log.Info(action)

	return nil
}

func (c *Cli) register() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин(действующий email): "
	login, err := c.prompt.PromptGetInput(loginPrompt, validator.ValidateEmail())
	if err != nil {
		return err
	}
	confirmationCode := c.codeGenerator.GetCode()
	confirmationMessage := "Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: " + confirmationCode
	err = c.mailer.SendMail([]string{login}, confirmationMessage)
	if err != nil {
		return err
	}
	codePrompt := prompt.PromptContent{}
	codePrompt.Label = "На указанный вами email отправлен код подтверждения. Введите код для продолжения регистрации: "
	_, err = c.prompt.PromptGetInput(codePrompt, validator.ValidateConfirmationCode(confirmationCode))
	if err != nil {
		return err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль(не менее 8 символов в разном регистре: буквы, цифры, спецсимволы.): ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validator.ValidatePassword())
	if err != nil {
		return err
	}
	c.log.Info(password)
	return nil
}
