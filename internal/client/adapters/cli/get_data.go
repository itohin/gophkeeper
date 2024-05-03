package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/validator"
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
		return "", err
	}
	switch s.SecretType {
	case entities.TypePassword:
		return c.showPassword(s)
	case entities.TypeText:
		return c.showText(s)
	case entities.TypeBinary:
		return c.showBinary(s)
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

func (c *Cli) showBinary(secret *entities.Secret) (string, error) {
	fmt.Println("Название: ", secret.Name)
	fmt.Println("Примечания: ", secret.Notes)
	action, err := c.prompt.PromptGetSelect(
		prompt.PromptContent{Label: "Выберите действие: "},
		[]prompt.SelectItem{
			{
				Label:  saveBinaryToDiskLabel,
				Action: saveBinaryToDisk,
			},
			{
				Label:  deleteDataLabel,
				Action: deleteData + "/" + secret.ID,
			},
			{
				Label:  comeBackLabel,
				Action: getData,
			},
		})
	if action != saveBinaryToDisk {
		return action, err
	}

	err = c.saveBinary(secret.Data.([]byte))
	if err != nil {
		return "", err
	}

	return getData, nil
}

func (c *Cli) saveBinary(data []byte) error {
	path, err := c.prompt.PromptGetInput(
		prompt.PromptContent{Label: "Введите путь для сохранение(/path/to/file.ext): "},
		validator.ValidateStringLength(3, 50),
	)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0755)
}

func (c *Cli) deleteData(id string) (string, error) {
	err := c.secrets.DeleteSecret(context.Background(), id)
	if err != nil {
		return "", err
	}
	return dataMenu, nil
}
