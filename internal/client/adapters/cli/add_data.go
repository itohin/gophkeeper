package cli

import (
	"context"
	"os"

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
	err = c.secrets.CreateSecret(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypeText,
			Data:       text,
		},
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
	err = c.secrets.CreateSecret(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypePassword,
			Data: &entities.Password{
				Login:    login,
				Password: password,
			},
		},
	)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}

func (c *Cli) addBinary() (string, error) {
	name, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите название: "},
		validator.ValidateStringLength(3, 25),
	)
	if err != nil {
		return "", err
	}
	path, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите путь к файлу: "},
		validator.ValidateStringLength(3, 500),
	)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(path)
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
	err = c.secrets.CreateSecret(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypeBinary,
			Data:       content,
		},
	)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}

func (c *Cli) addCard() (string, error) {
	name, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите название: "},
		validator.ValidateStringLength(3, 25),
	)
	if err != nil {
		return "", err
	}
	number, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите номер карты: "},
		validator.ValidateStringLength(13, 25),
	)
	if err != nil {
		return "", err
	}
	expiration, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите срок окончания действия карты в формате ММ/ГГГГ: "},
		validator.ValidateCardExpiration(),
	)
	if err != nil {
		return "", err
	}
	code, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите cvc/cvv код карты: "},
		validator.ValidateStringLength(0, 3),
	)
	if err != nil {
		return "", err
	}
	pin, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите pin код карты: "},
		validator.ValidateStringLength(0, 25),
	)
	if err != nil {
		return "", err
	}
	ownerName, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите имя владельца карты: "},
		validator.ValidateStringLength(0, 25),
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
	err = c.secrets.CreateSecret(
		context.Background(),
		&entities.Secret{
			Name:       name,
			Notes:      notes,
			SecretType: entities.TypeCard,
			Data: &entities.Card{
				Number:     number,
				Expiration: expiration,
				Code:       code,
				Pin:        pin,
				OwnerName:  ownerName,
			},
		},
	)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}
