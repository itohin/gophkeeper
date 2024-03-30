package cli

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/validator"
)

func (c *Cli) addText() (string, error) {
	name, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите название: "},
		validator.ValidateStringLength(3, 25),
	)
	if err != nil {
		return "", err
	}
	text, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите текст: "},
		validator.ValidateStringLength(3, 500),
	)
	if err != nil {
		return "", err
	}
	notes, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите примечания: "},
		validator.ValidateStringLength(0, 500),
	)
	if err != nil {
		return "", err
	}
	err = c.secrets.CreateText(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypeText,
		},
		text,
	)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}

func (c *Cli) addPassword() (string, error) {
	name, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите название: "},
		validator.ValidateStringLength(3, 25),
	)
	if err != nil {
		return "", err
	}
	login, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите логин: "},
		validator.ValidateStringLength(3, 30),
	)
	if err != nil {
		return "", err
	}
	password, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите пароль: ", Mask: 42},
		validator.ValidateStringLength(5, 30),
	)
	if err != nil {
		return "", err
	}
	notes, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите примечания: "},
		validator.ValidateStringLength(0, 500),
	)
	if err != nil {
		return "", err
	}
	err = c.secrets.CreatePassword(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypePassword,
		},
		&entities.Password{
			Login:    login,
			Password: password,
		},
	)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}
