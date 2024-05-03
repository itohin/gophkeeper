package cli

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCli_addText(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	secrets := mocks.NewMockSecrets(ctrl)

	c := &Cli{
		prompt:  prompter,
		secrets: secrets,
	}

	namePrompt := prompt.PromptContent{}
	namePrompt.Label = "Введите название: "

	textPrompt := prompt.PromptContent{}
	textPrompt.Label = "Введите текст: "

	notesPrompt := prompt.PromptContent{}
	notesPrompt.Label = "Введите примечания: "

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "name error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"textPrompt":   0,
				"notesPrompt":  0,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   errors.New("name error"),
				"textPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "text error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"textPrompt":   1,
				"notesPrompt":  0,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"textPrompt":   errors.New("text error"),
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "notes error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"textPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"textPrompt":   nil,
				"notesPrompt":  errors.New("notes error"),
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "create text error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"textPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 1,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"textPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": errors.New("create error"),
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"textPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 1,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"textPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "dataMenu",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			prompter.EXPECT().PromptGetInput(namePrompt, gomock.Any()).Return("First", tt.errors["namePrompt"]).Times(tt.mockTimes["namePrompt"])
			prompter.EXPECT().PromptGetInput(textPrompt, gomock.Any()).Return("Lorem ipsum...", tt.errors["textPrompt"]).Times(tt.mockTimes["textPrompt"])
			prompter.EXPECT().PromptGetInput(notesPrompt, gomock.Any()).Return("first notes", tt.errors["notesPrompt"]).Times(tt.mockTimes["notesPrompt"])
			secrets.EXPECT().CreateSecret(gomock.Any(), &entities.Secret{
				Name:       "First",
				Notes:      "first notes",
				SecretType: 1,
				Data:       "Lorem ipsum...",
			}).Return(tt.errors["createSecret"]).Times(tt.mockTimes["createSecret"])

			_, err := c.addText()
			tt.wantErr(t, err, fmt.Sprintf("AddText()"))
		})
	}
}

func TestCli_addPassword(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	secrets := mocks.NewMockSecrets(ctrl)

	c := &Cli{
		prompt:  prompter,
		secrets: secrets,
	}

	namePrompt := prompt.PromptContent{}
	namePrompt.Label = "Введите название: "

	loginPrompt := prompt.PromptContent{}
	loginPrompt.Label = "Введите логин: "

	passwordPrompt := prompt.PromptContent{
		Label: "Введите пароль: ",
		Mask:  42,
	}

	notesPrompt := prompt.PromptContent{}
	notesPrompt.Label = "Введите примечания: "

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "name error",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    0,
				"passwordPrompt": 0,
				"notesPrompt":    0,
				"createSecret":   0,
			},
			errors: map[string]error{
				"namePrompt":     errors.New("name error"),
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"notesPrompt":    nil,
				"createSecret":   nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "login error",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    1,
				"passwordPrompt": 0,
				"notesPrompt":    0,
				"createSecret":   0,
			},
			errors: map[string]error{
				"namePrompt":     nil,
				"loginPrompt":    errors.New("login error"),
				"passwordPrompt": nil,
				"notesPrompt":    nil,
				"createSecret":   nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "password error",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"notesPrompt":    0,
				"createSecret":   0,
			},
			errors: map[string]error{
				"namePrompt":     nil,
				"loginPrompt":    nil,
				"passwordPrompt": errors.New("password error"),
				"notesPrompt":    nil,
				"createSecret":   nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "notes error",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"notesPrompt":    1,
				"createSecret":   0,
			},
			errors: map[string]error{
				"namePrompt":     nil,
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"notesPrompt":    errors.New("notes error"),
				"createSecret":   nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "create text error",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"notesPrompt":    1,
				"createSecret":   1,
			},
			errors: map[string]error{
				"namePrompt":     nil,
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"notesPrompt":    nil,
				"createSecret":   errors.New("create error"),
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"namePrompt":     1,
				"loginPrompt":    1,
				"passwordPrompt": 1,
				"notesPrompt":    1,
				"createSecret":   1,
			},
			errors: map[string]error{
				"namePrompt":     nil,
				"loginPrompt":    nil,
				"passwordPrompt": nil,
				"notesPrompt":    nil,
				"createSecret":   nil,
			},
			wantAction: "dataMenu",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			prompter.EXPECT().PromptGetInput(namePrompt, gomock.Any()).Return("First pass", tt.errors["namePrompt"]).Times(tt.mockTimes["namePrompt"])
			prompter.EXPECT().PromptGetInput(loginPrompt, gomock.Any()).Return("aaa@zzz.com", tt.errors["loginPrompt"]).Times(tt.mockTimes["loginPrompt"])
			prompter.EXPECT().PromptGetInput(passwordPrompt, gomock.Any()).Return("pass@Word1", tt.errors["passwordPrompt"]).Times(tt.mockTimes["passwordPrompt"])
			prompter.EXPECT().PromptGetInput(notesPrompt, gomock.Any()).Return("http://aaa.zzz", tt.errors["notesPrompt"]).Times(tt.mockTimes["notesPrompt"])
			secrets.EXPECT().CreateSecret(gomock.Any(), &entities.Secret{
				Name:       "First pass",
				Notes:      "http://aaa.zzz",
				SecretType: 2,
				Data: &entities.Password{
					Login:    "aaa@zzz.com",
					Password: "pass@Word1",
				},
			}).Return(tt.errors["createSecret"]).Times(tt.mockTimes["createSecret"])

			_, err := c.addPassword()
			tt.wantErr(t, err, fmt.Sprintf("AddPassword()"))
		})
	}
}

func TestCli_addBinary(t *testing.T) {

	file, _ := os.Create("test_data.txt")
	file.Write([]byte("Lorem ipsum..."))
	defer os.Remove("test_data.txt")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	secrets := mocks.NewMockSecrets(ctrl)

	c := &Cli{
		prompt:  prompter,
		secrets: secrets,
	}

	namePrompt := prompt.PromptContent{}
	namePrompt.Label = "Введите название: "

	pathPrompt := prompt.PromptContent{}
	pathPrompt.Label = "Введите путь к файлу: "

	notesPrompt := prompt.PromptContent{}
	notesPrompt.Label = "Введите примечания: "

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "name error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"pathPrompt":   0,
				"notesPrompt":  0,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   errors.New("name error"),
				"pathPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "path error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"pathPrompt":   1,
				"notesPrompt":  0,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"pathPrompt":   errors.New("path error"),
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "notes error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"pathPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 0,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"pathPrompt":   nil,
				"notesPrompt":  errors.New("notes error"),
				"createSecret": nil,
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "create binary error",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"pathPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 1,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"pathPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": errors.New("create error"),
			},
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "success",
			mockTimes: map[string]int{
				"namePrompt":   1,
				"pathPrompt":   1,
				"notesPrompt":  1,
				"createSecret": 1,
			},
			errors: map[string]error{
				"namePrompt":   nil,
				"pathPrompt":   nil,
				"notesPrompt":  nil,
				"createSecret": nil,
			},
			wantAction: "dataMenu",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompter.EXPECT().PromptGetInput(namePrompt, gomock.Any()).Return("Binary", tt.errors["namePrompt"]).Times(tt.mockTimes["namePrompt"])
			prompter.EXPECT().PromptGetInput(pathPrompt, gomock.Any()).Return("test_data.txt", tt.errors["pathPrompt"]).Times(tt.mockTimes["pathPrompt"])
			prompter.EXPECT().PromptGetInput(notesPrompt, gomock.Any()).Return("binary notes", tt.errors["notesPrompt"]).Times(tt.mockTimes["notesPrompt"])
			secrets.EXPECT().CreateSecret(gomock.Any(), &entities.Secret{
				Name:       "Binary",
				Notes:      "binary notes",
				SecretType: 3,
				Data:       []byte("Lorem ipsum..."),
			}).Return(tt.errors["createSecret"]).Times(tt.mockTimes["createSecret"])

			_, err := c.addBinary()
			tt.wantErr(t, err, fmt.Sprintf("AddBinary()"))
		})
	}
}
