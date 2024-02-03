package cli

import (
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
)

func (c *Cli) dataMenu() (string, error) {
	return c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите действие: "},
		[]prompt.SelectItem{
			{
				Label:  addDataLabel,
				Action: addData,
			},
			{
				Label:  getDataLabel,
				Action: getData,
			},
		})
}

func (c *Cli) addData() (string, error) {
	c.log.Info("add data")
	return "", nil
}

func (c *Cli) getData() (string, error) {
	c.log.Info("get data")
	return "", nil
}
