package cli

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/pkg/validator"
)

const (
	register = "Регистрация"
	login    = "Вход"
)

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
	login, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите логин: "},
		validator.ValidateEmail(),
	)
	if err != nil {
		return err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validator.ValidatePassword())
	if err != nil {
		return err
	}

	err = c.auth.Login(context.Background(), login, password)
	if err != nil {
		return err
	}
	return c.dataMenu()
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

	err = c.auth.Register(context.Background(), login, password)
	if err != nil {
		return err
	}

	return c.dataMenu()
}
