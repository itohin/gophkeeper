package cli

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/adapters/cli/router"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCli_Call(t *testing.T) {

	tests := []struct {
		name       string
		routes     map[string]router.Command
		action     string
		wantResult string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "empty action error",
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:       "no routes for action error",
			action:     "action",
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "action type error",
			action: "action",
			routes: map[string]router.Command{
				"action": 111,
			},
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "not enough parameters error",
			action: "action",
			routes: map[string]router.Command{
				"action": func(id string) (string, error) { return "", nil },
			},
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "more parameters than needed error",
			action: "action/id",
			routes: map[string]router.Command{
				"action": func() (string, error) { return "", nil },
			},
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "not enough parameters returned error",
			action: "action",
			routes: map[string]router.Command{
				"action": func() error { return nil },
			},
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "action error success",
			action: "action",
			routes: map[string]router.Command{
				"action": func() (string, error) { return "", fmt.Errorf("any error") },
			},
			wantResult: "",
			wantErr:    assert.Error,
		},
		{
			name:   "action success",
			action: "action",
			routes: map[string]router.Command{
				"action": func() (string, error) { return "hello", nil },
			},
			wantResult: "hello",
			wantErr:    assert.NoError,
		},
		{
			name:   "action with param success",
			action: "action/id",
			routes: map[string]router.Command{
				"action": func(id string) (string, error) { return id, nil },
			},
			wantResult: "id",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cli{router: router.NewRouter(tt.routes)}

			gotResult, err := c.Call(tt.action)
			if !tt.wantErr(t, err, fmt.Sprintf("Call(%v)", tt.action)) {
				return
			}
			assert.Equalf(t, tt.wantResult, gotResult, "Call(%v)", tt.action)
		})
	}
}
