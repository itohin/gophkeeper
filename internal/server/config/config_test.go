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
				&DB{
					DSN:            "postgres://postgres:postgres@localhost:5432/gophkeeper",
					MigrationsPath: "internal/server/infrastructure/migrations",
				},
				&JWT{
					Signature:  "secret",
					AccessTTL:  60 * time.Second,
					RefreshTTL: 360 * time.Second,
				},
				&WebSocket{
					Address: ":7777",
				},
				&SSL{
					CertPath: "ca.crt",
					KeyPath:  "ca.key",
				},
				&Mail{
					Login:    "from@gmail.com",
					Password: "",
					Host:     "localhost",
					Port:     "1025",
				},
				&GRPC{
					Address: ":3200",
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
					"DB_DSN":             "postgres://postgres:postgres@envhost:5432/gophkeeper",
					"DB_MIGRATIONS_PATH": "env/migrations",
					"JWT_SIGNATURE":      "envsecret",
					"JWT_ACCESS_TTL":     "15s",
					"JWT_REFRESH_TTL":    "35s",
					"WEBSOCKET_ADDRESS":  ":9999",
					"SSL_CERT_PATH":      "env.crt",
					"SSL_KEY_PATH":       "env.key",
					"MAIL_LOGIN":         "env@mail.ru",
					"MAIL_PASSWORD":      "env",
					"MAIL_HOST":          "envhost",
					"MAIL_PORT":          "1026",
					"GRPC_ADDRESS":       ":3400",
				},
			},
			want: &AppConfig{
				&DB{
					DSN:            "postgres://postgres:postgres@envhost:5432/gophkeeper",
					MigrationsPath: "env/migrations",
				},
				&JWT{
					Signature:  "envsecret",
					AccessTTL:  15 * time.Second,
					RefreshTTL: 35 * time.Second,
				},
				&WebSocket{
					Address: ":9999",
				},
				&SSL{
					CertPath: "env.crt",
					KeyPath:  "env.key",
				},
				&Mail{
					Login:    "env@mail.ru",
					Password: "env",
					Host:     "envhost",
					Port:     "1026",
				},
				&GRPC{
					Address: ":3400",
				},
			},
		},
		{
			name: "flags",
			arrange: arrange{
				args: []string{
					"cmd",
					"--db-dsn=postgres://postgres:postgres@flaghost:5432/gophkeeper",
					"--db-mig-path=flag/migrations",
					"--jwt-sig=flagsecret",
					"--jwt-attl=10s",
					"--jwt-rttl=30s",
					"--ws-addr=:8888",
					"--ssl-cert=flag.crt",
					"--ssl-key=flag.key",
					"--mail-login=flag@mail.ru",
					"--mail-pass=flag",
					"--mail-host=flaghost",
					"--mail-port=1027",
					"--grpc-addr=:3300",
				},
				env: map[string]string{
					"DB_DSN":             "postgres://postgres:postgres@envhost:5432/gophkeeper",
					"DB_MIGRATIONS_PATH": "env/migrations",
					"JWT_SIGNATURE":      "envsecret",
					"JWT_ACCESS_TTL":     "15s",
					"JWT_REFRESH_TTL":    "35s",
					"WEBSOCKET_ADDRESS":  ":9999",
					"SSL_CERT_PATH":      "env.crt",
					"SSL_KEY_PATH":       "env.key",
					"MAIL_LOGIN":         "env@mail.ru",
					"MAIL_PASSWORD":      "env",
					"MAIL_HOST":          "envhost",
					"MAIL_PORT":          "1026",
					"GRPC_ADDRESS":       ":3400",
				},
			},
			want: &AppConfig{
				&DB{
					DSN:            "postgres://postgres:postgres@flaghost:5432/gophkeeper",
					MigrationsPath: "flag/migrations",
				},
				&JWT{
					Signature:  "flagsecret",
					AccessTTL:  10 * time.Second,
					RefreshTTL: 30 * time.Second,
				},
				&WebSocket{
					Address: ":8888",
				},
				&SSL{
					CertPath: "flag.crt",
					KeyPath:  "flag.key",
				},
				&Mail{
					Login:    "flag@mail.ru",
					Password: "flag",
					Host:     "flaghost",
					Port:     "1027",
				},
				&GRPC{
					Address: ":3300",
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
