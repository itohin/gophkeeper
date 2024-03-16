package cli

import "github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"

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
			{
				Label:  logoutLabel,
				Action: logout,
			},
		})
}

func (c *Cli) addData() (string, error) {
	return c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите тип данных: "},
		[]prompt.SelectItem{
			{
				Label:  addTextLabel,
				Action: addText,
			},
			{
				Label:  addPasswordLabel,
				Action: addPassword,
			},
			{
				Label:  comeBackLabel,
				Action: dataMenu,
			},
		})
}
