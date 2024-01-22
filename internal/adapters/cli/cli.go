package cli

import "C"
import (
	"errors"
	"fmt"
	"github.com/itohin/gophkeeper/internal/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/infrastructure/logger"
	"github.com/itohin/gophkeeper/internal/infrastructure/mailer"
	"github.com/itohin/gophkeeper/internal/infrastructure/validator"
	"math/rand"
	"strconv"
	"time"
)

const (
	register = "Регистрация"
	login    = "Вход"
	addData  = "Сохранить данные"
	getData  = "Получить данные"
)

type Cli struct {
	log    logger.Logger
	prompt *prompt.Prompt
	mailer mailer.Mailer
}

func NewCli(logger logger.Logger, prompt *prompt.Prompt, mailer mailer.Mailer) *Cli {
	return &Cli{
		log:    logger,
		prompt: prompt,
		mailer: mailer,
	}
}

func (c *Cli) Auth() error {
	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выполните вход или зарегистрируйтесь: "

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{login, register})
	if err != nil {
		return err
	}
	switch action {
	case login:
		return c.Login()
	case register:
		return c.Register()
	default:
		return fmt.Errorf("unknown action %s", action)
	}
}

func (c *Cli) Login() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин: "
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("error")
		}
		return nil
	}
	login, err := c.prompt.PromptGetInput(loginPrompt, validate)
	if err != nil {
		return err
	}
	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}
	password, err := c.prompt.PromptGetInput(passwordPrompt, validate)
	if err != nil {
		return err
	}

	c.log.Info(login, password)
	//check credentials

	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выберите действие: "

	action, err := c.prompt.PromptGetSelect(menuPrompt, []string{addData, getData})
	if err != nil {
		return err
	}
	c.log.Info(action)

	return nil
}

func (c *Cli) Register() error {
	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин(действующий email): "
	login, err := c.prompt.PromptGetInput(loginPrompt, validator.ValidateEmail())
	if err != nil {
		return err
	}
	confirmationCode := getConfirmationCode()
	confirmationMessage := "Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: " + confirmationCode
	err = c.mailer.SendMail([]string{login}, confirmationMessage)
	if err != nil {
		return err
	}
	codePrompt := prompt.PromptContent{}
	codePrompt.Label = "На указанный вами email отправлен код подтверждения. Введите код для продолжения регистрации: "
	code, err := c.prompt.PromptGetInput(codePrompt, validator.ValidateConfirmationCode(confirmationCode))
	if err != nil {
		return err
	}
	c.log.Info(code)
	return nil
}

func getConfirmationCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	min := 1001
	max := 9999
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}
