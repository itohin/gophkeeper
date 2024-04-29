package cli

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/prompt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/mocks"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCli_showData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	secrets := mocks.NewMockSecrets(ctrl)

	c := &Cli{
		prompt:  prompter,
		secrets: secrets,
	}

	secretsMap := map[string]*entities.Secret{
		"text": &entities.Secret{
			ID:         "text",
			Name:       "First",
			SecretType: 1,
			Data:       "Lorem ipsum...",
			Notes:      "first notes",
		},
		"password": &entities.Secret{
			ID:         "password",
			Name:       "Password",
			SecretType: 2,
			Data:       &entities.Password{Login: "aaa@zzz.com", Password: "pass@Word1"},
			Notes:      "http://aaa.zzz",
		},
		"binary": &entities.Secret{
			ID:         "binary",
			Name:       "Binary",
			SecretType: 3,
			Data:       []byte("Lorem ipsum..."),
			Notes:      "binary notes",
		},
		"unknown_type": &entities.Secret{
			ID:         "unknown_type",
			Name:       "Unknown",
			SecretType: 4,
			Data:       "Unknown type",
			Notes:      "...",
		},
	}

	menuPrompt := prompt.PromptContent{}
	menuPrompt.Label = "Выберите действие: "

	tests := []struct {
		name       string
		mockTimes  map[string]int
		errors     map[string]error
		id         string
		wantAction string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "get secret error",
			mockTimes: map[string]int{
				"promptSelect": 0,
			},
			errors: map[string]error{
				"getSecret":    errors.New("get secret error"),
				"promptSelect": nil,
			},
			id:         "text",
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "unknown secret type error",
			mockTimes: map[string]int{
				"promptSelect": 0,
			},
			errors: map[string]error{
				"getSecret":    nil,
				"promptSelect": nil,
			},
			id:         "unknown_type",
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "unknown secret type error",
			mockTimes: map[string]int{
				"promptSelect": 0,
			},
			errors: map[string]error{
				"getSecret":    nil,
				"promptSelect": nil,
			},
			id:         "unknown_type",
			wantAction: "",
			wantErr:    assert.Error,
		},
		{
			name: "delete action success",
			mockTimes: map[string]int{
				"promptSelect": 1,
			},
			errors: map[string]error{
				"getSecret":    nil,
				"promptSelect": nil,
			},
			id:         "text",
			wantAction: "deleteData/text",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			menuItems := []prompt.SelectItem{
				{
					Label:  "Удалить данные",
					Action: "deleteData/" + tt.id,
				},
				{
					Label:  "Вернуться назад",
					Action: "getData",
				},
			}

			secrets.EXPECT().GetSecret(gomock.Any(), tt.id).Return(secretsMap[tt.id], tt.errors["getSecret"]).Times(1)
			prompter.EXPECT().PromptGetSelect(menuPrompt, menuItems).Return(tt.wantAction, tt.errors["promptSelect"]).Times(tt.mockTimes["promptSelect"])

			_, err := c.showData(tt.id)
			tt.wantErr(t, err, fmt.Sprintf("ShowData()"))
		})
	}
}

func TestCli_saveBinary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	prompter := mocks.NewMockPrompter(ctrl)
	c := &Cli{prompt: prompter}
	data := []byte("Lorem ipsum...")

	tests := []struct {
		name    string
		errors  map[string]error
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Prompt error",
			errors: map[string]error{
				"promptError": errors.New("prompt error"),
			},
			wantErr: assert.Error,
		},
		{
			name: "success",
			errors: map[string]error{
				"promptError": nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			p := prompt.PromptContent{}
			p.Label = "Введите путь для сохранение(/path/to/file.ext): "
			prompter.EXPECT().PromptGetInput(p, gomock.Any()).Return("data.txt", tt.errors["promptError"]).Times(1)

			tt.wantErr(t, c.saveBinary(data), fmt.Sprintf("ShowData()"))

			if tt.errors["promptError"] == nil {
				fileData, _ := os.ReadFile("data.txt")
				os.Remove("data.txt")
				assert.Equal(t, fileData, []byte("Lorem ipsum..."))
			}
		})
	}
}

func TestCli_deleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	secrets := mocks.NewMockSecrets(ctrl)
	c := &Cli{secrets: secrets}

	tests := []struct {
		name    string
		error   error
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "delete error",
			error:   errors.New("delete error"),
			want:    "",
			wantErr: assert.Error,
		},
		{
			name:    "success",
			error:   nil,
			want:    "dataMenu",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secrets.EXPECT().DeleteSecret(gomock.Any(), "uuid-001").Return(tt.error).Times(1)

			got, err := c.deleteData("uuid-001")
			if !tt.wantErr(t, err, fmt.Sprintf("deleteData")) {
				return
			}
			assert.Equalf(t, tt.want, got, "deleteData")
		})
	}
}
