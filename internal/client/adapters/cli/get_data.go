package cli

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

func (c *Cli) getData() (string, error) {
	secrets, err := c.secrets.GetSecrets(context.Background())
	if err != nil {
		return "", err
	}
	menu := make([]prompt.SelectItem, 0, len(secrets))
	for _, secret := range secrets {
		p := prompt.SelectItem{
			Label:  secret.Name + " (" + secret.GetLabel() + ")",
			Action: showData + "/" + secret.ID,
		}
		menu = append(menu, p)
	}
	menu = append(menu, prompt.SelectItem{
		Label:  comeBackLabel,
		Action: dataMenu,
	})

	listPrompt := prompt.PromptContent{}
	listPrompt.Label = "Выберите запись: "
	return c.prompt.PromptGetSelect(listPrompt, menu)

}

func (c *Cli) showData(id string) (string, error) {
	s, err := c.secrets.GetSecret(context.Background(), id)
	if err != nil {
		return "", nil
	}
	switch s.SecretType {
	case entities.TypePassword:
		return c.showPassword(s)
	case entities.TypeText:
		return c.showText(s)
	default:
		return "", fmt.Errorf("unknown secret type %v", s.SecretType)
	}
}

func (c *Cli) showPassword(secret *entities.Secret) (string, error) {
	data := secret.Data.(*entities.Password)
	fmt.Println("Название: ", secret.Name)
	fmt.Println("Логин: ", data.Login)
	fmt.Println("Пароль: ", data.Password)
	fmt.Println("Примечания: ", secret.Notes)
	return c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите действие: "},
		[]prompt.SelectItem{
			{
				Label:  deleteDataLabel,
				Action: deleteData + "/" + secret.ID,
			},
			{
				Label:  comeBackLabel,
				Action: getData,
			},
		})
}

func (c *Cli) showText(secret *entities.Secret) (string, error) {
	fmt.Println("Название: ", secret.Name)
	fmt.Println("Текст: ", secret.Data.(string))
	fmt.Println("Примечания: ", secret.Notes)
	return c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите действие: "},
		[]prompt.SelectItem{
			{
				Label:  deleteDataLabel,
				Action: deleteData + "/" + secret.ID,
			},
			{
				Label:  comeBackLabel,
				Action: getData,
			},
		})
}

func (c *Cli) deleteData(id string) (string, error) {
	err := c.secrets.DeleteSecret(context.Background(), id)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}
