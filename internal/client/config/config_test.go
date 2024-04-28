package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadConfig(t *testing.T) {
	type arrange struct {
		args []string
		env  map[string]string
	}
	tests := []struct {
		name    string
		arrange arrange
		want    *AppConfig
	}{
		{
			name: "default",
			arrange: arrange{
				args: []string{
					"cmd",
				},
				env: map[string]string{},
			},
			want: &AppConfig{
				&JWT{
					Signature:  "secret",
					AccessTTL:  60 * time.Second,
					RefreshTTL: 360 * time.Second,
				},
				&WebSocket{
					ServerAddress:     ":7777",
					ConnectionTimeout: 100 * time.Millisecond,
				},
				&GRPC{
					ServerAddress: ":3200",
				},
			},
		},
		{
			name: "env",
			arrange: arrange{
				args: []string{
					"cmd",
				},
				env: map[string]string{
					"JWT_SIGNATURE":                "envsecret",
					"JWT_ACCESS_TTL":               "15s",
					"JWT_REFRESH_TTL":              "35s",
					"WEBSOCKET_ADDRESS":            ":9999",
					"WEBSOCKET_CONNECTION_TIMEOUT": "300ms",
					"GRPC_ADDRESS":                 ":3400",
				},
			},
			want: &AppConfig{
				&JWT{
					Signature:  "envsecret",
					AccessTTL:  15 * time.Second,
					RefreshTTL: 35 * time.Second,
				},
				&WebSocket{
					ServerAddress:     ":9999",
					ConnectionTimeout: 300 * time.Millisecond,
				},
				&GRPC{
					ServerAddress: ":3400",
				},
			},
		},
		{
			name: "flags",
			arrange: arrange{
				args: []string{
					"cmd",
					"--jwt-sig=flagsecret",
					"--jwt-attl=10s",
					"--jwt-rttl=30s",
					"--ws-addr=:8888",
					"--ws-ttl=200ms",
					"--grpc-addr=:3300",
				},
				env: map[string]string{},
			},
			want: &AppConfig{
				&JWT{
					Signature:  "flagsecret",
					AccessTTL:  10 * time.Second,
					RefreshTTL: 30 * time.Second,
				},
				&WebSocket{
					ServerAddress:     ":8888",
					ConnectionTimeout: 200 * time.Millisecond,
				},
				&GRPC{
					ServerAddress: ":3300",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.arrange.args
			for k, v := range tt.arrange.env {
				os.Setenv(k, v)
			}

			if got := ReadConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got.JWT, tt.want.JWT)
			}
		})
	}
}
