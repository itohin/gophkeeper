package cli

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCli_Auth(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)

	c := &Cli{prompt: prompter}

	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выполните вход или зарегистрируйтесь: "
	selectItems := []prompt.SelectItem{
		{
			Label:  loginLabel,
			Action: login,
		},
		{
			Label:  registerLabel,
			Action: register,
		},
	}

	tests := []struct {
		name       string
		err        error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "unknown action",
			err:        errors.New("unknown action"),
			wantAction: "unknown",
			wantErr:    assert.Error,
		},
		{
			name:       "login",
			err:        nil,
			wantAction: "login",
			wantErr:    assert.NoError,
		},
		{
			name:       "register",
			err:        nil,
			wantAction: "register",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompter.EXPECT().PromptGetSelect(menuPrompt, selectItems).Return(tt.wantAction, tt.err).Times(1)

			action, err := c.authMenu()

			assert.Equal(t, action, tt.wantAction)

			tt.wantErr(t, err, fmt.Sprintf("authMenu()"))
		})
	}
}

func TestCli_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	mailer := mocks.NewMockMailer(ctrl)
	log := mocks.NewMockLogger(ctrl)
	codegen := mocks.NewMockGenerator(ctrl)
	auth := mocks.NewMockAuth(ctrl)

	log.EXPECT().Info(gomock.Any()).AnyTimes()

	c := &Cli{
		log:           log,
		prompt:        prompter,
		mailer:        mailer,
		codeGenerator: codegen,
		auth:          auth,
	}

	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин(действующий email): "

	codePrompt := prompt.PromptContent{}
	codePrompt.Label = "На указанный вами email отправлен код подтверждения. Введите код для продолжения регистрации: "

	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль(не менее 8 символов в разном регистре: буквы, цифры, спецсимволы.): ",
		Mask:  42,
	}

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "login error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        0,
				"sendMail":       0,
				"codePrompt":     0,
				"passwordPrompt": 0,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    errors.New("login error"),
				"sendMail":       nil,
				"codePrompt":     nil,
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "mail error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        1,
				"sendMail":       1,
				"codePrompt":     0,
				"passwordPrompt": 0,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"sendMail":       errors.New("mail error"),
				"codePrompt":     nil,
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "code error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        1,
				"sendMail":       1,
				"codePrompt":     1,
				"passwordPrompt": 0,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"sendMail":       nil,
				"codePrompt":     errors.New("code error"),
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "password error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        1,
				"sendMail":       1,
				"codePrompt":     1,
				"passwordPrompt": 1,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"sendMail":       nil,
				"codePrompt":     nil,
				"passwordPrompt": errors.New("password error"),
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "auth register error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        1,
				"sendMail":       1,
				"codePrompt":     1,
				"passwordPrompt": 1,
				"auth":           1,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"sendMail":       nil,
				"codePrompt":     nil,
				"passwordPrompt": nil,
				"auth":           errors.New("register error"),
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"codeGen":        1,
				"sendMail":       1,
				"codePrompt":     1,
				"passwordPrompt": 1,
				"auth":           1,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"sendMail":       nil,
				"codePrompt":     nil,
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "dataMenu",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			prompter.EXPECT().PromptGetInput(loginPrompt, gomock.Any()).Return("a@a.com", tt.errors["loginPrompt"]).Times(tt.mockTimes["loginPrompt"])
			codegen.EXPECT().GetCode().Return("1111").Times(tt.mockTimes["codeGen"])
			mailer.EXPECT().SendMail([]string{"a@a.com"}, "Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: 1111").Return(tt.errors["sendMail"]).Times(tt.mockTimes["sendMail"])
			prompter.EXPECT().PromptGetInput(codePrompt, gomock.Any()).Return("1111", tt.errors["codePrompt"]).Times(tt.mockTimes["codePrompt"])
			prompter.EXPECT().PromptGetInput(passwordPrompt, gomock.Any()).Return("tesT@pass1word", tt.errors["passwordPrompt"]).Times(tt.mockTimes["passwordPrompt"])
			auth.EXPECT().Register(gomock.Any(), "a@a.com", "tesT@pass1word").Return(tt.errors["auth"]).Times(tt.mockTimes["auth"])

			_, err := c.register()
			tt.wantErr(t, err, fmt.Sprintf("Register()"))
		})
	}
}

func TestCli_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	log := mocks.NewMockLogger(ctrl)
	auth := mocks.NewMockAuth(ctrl)

	log.EXPECT().Info(gomock.Any()).AnyTimes()

	c := &Cli{
		log:    log,
		prompt: prompter,
		auth:   auth,
	}

	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин: "

	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "login error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"passwordPrompt": 0,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    errors.New("login error"),
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "password error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"auth":           0,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"passwordPrompt": errors.New("password error"),
				"auth":           nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "auth login error",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"auth":           1,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"auth":           errors.New("login error"),
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"auth":           1,
			},
			errors: map[string]error{
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"auth":           nil,
			},
			wantAction: "dataMenu",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			prompter.EXPECT().PromptGetInput(loginPrompt, gomock.Any()).Return("a@a.com", tt.errors["loginPrompt"]).Times(tt.mockTimes["loginPrompt"])
			prompter.EXPECT().PromptGetInput(passwordPrompt, gomock.Any()).Return("tesT@pass1word", tt.errors["passwordPrompt"]).Times(tt.mockTimes["passwordPrompt"])
			auth.EXPECT().Login(gomock.Any(), "a@a.com", "tesT@pass1word").Return(tt.errors["auth"]).Times(tt.mockTimes["auth"])

			action, err := c.login()
			assert.Equal(t, action, tt.wantAction)
			tt.wantErr(t, err, fmt.Sprintf("Login()"))
		})
	}
}
