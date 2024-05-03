package cli

import (
	"context"

	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/pkg/validator"
)

func (c *Cli) authMenu() (string, error) {
	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выполните вход или зарегистрируйтесь: "

	return c.prompt.PromptGetSelect(menuPrompt, []prompt.SelectItem{
		{
			Label:  loginLabel,
			Action: login,
		},
		{
			Label:  registerLabel,
			Action: register,
		},
		{
			Label:  verifyLabel,
			Action: verify,
		},
		{
			Label:  logoutLabel,
			Action: logout,
		},
	})
}

func (c *Cli) login() (string, error) {
	login, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите логин: "},
		validator.ValidateEmail(),
	)
	if err != nil {
		return "", err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validator.ValidatePassword())
	if err != nil {
		return "", err
	}

	err = c.auth.Login(context.Background(), login, password)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}

func (c *Cli) register() (string, error) {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин(действующий email): "
	login, err := c.prompt.PromptGetInput(loginPrompt, validator.ValidateEmail())
	if err != nil {
		return "", err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль(не менее 8 символов в разном регистре: буквы, цифры, спецсимволы.): ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validator.ValidatePassword())
	if err != nil {
		return "", err
	}

	err = c.auth.Register(context.Background(), login, password)
	if err != nil {
		return "", err
	}
	codePrompt := prompt.PromptContent{}
	codePrompt.Label = "На указанный вами email отправлен код подтверждения. Введите код для продолжения регистрации: "
	otp, err := c.prompt.PromptGetInput(codePrompt, validator.ValidateConfirmationCode())
	if err != nil {
		return "", err
	}

	err = c.auth.Verify(context.Background(), login, otp)
	if err != nil {
		return "", err
	}

	return dataMenu, nil
}

func (c *Cli) verify() (string, error) {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите указанный при регистрации email: "
	login, err := c.prompt.PromptGetInput(loginPrompt, validator.ValidateEmail())
	if err != nil {
		return "", err
	}
	codePrompt := prompt.PromptContent{}
	codePrompt.Label = "Введите код подтверждения, полученный по email, для продолжения регистрации: "
	otp, err := c.prompt.PromptGetInput(codePrompt, validator.ValidateConfirmationCode())
	if err != nil {
		return "", err
	}

	err = c.auth.Verify(context.Background(), login, otp)
	if err != nil {
		return "", err
	}

	return dataMenu, nil
}

func (c *Cli) logout() (string, error) {
	return "", c.auth.Logout(context.Background())
}
