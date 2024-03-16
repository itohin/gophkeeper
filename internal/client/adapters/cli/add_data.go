package cli

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/pkg/validator"
)

func (c *Cli) addText() (string, error) {
	dataLabel, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите название: "},
		validator.ValidateStringLength(3, 25),
	)
	if err != nil {
		return "", err
	}
	data, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите текст: "},
		validator.ValidateStringLength(3, 500),
	)
	if err != nil {
		return "", err
	}
	fmt.Println(dataLabel, data)
	return "", nil
}

func (c *Cli) addPassword() (string, error) {
	c.log.Info("add password")
	return "", nil
}
