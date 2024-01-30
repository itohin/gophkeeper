package cli

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
)

const (
	addData = "Сохранить данные"
	getData = "Получить данные"
)

func (c *Cli) dataMenu() error {
	action, err := c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите действие: "},
		[]string{addData, getData},
	)
	if err != nil {
		return err
	}

	switch action {
	case addData:
		return c.addData()
	case getData:
		return c.getData()
	default:
		return fmt.Errorf("unknown action %s", action)
	}
}

func (c *Cli) addData() error {
	c.log.Info("add data")
	return nil
}

func (c *Cli) getData() error {
	c.log.Info("get data")
	return nil
}
