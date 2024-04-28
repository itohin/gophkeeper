package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	DBDsn            = "DbDsn"
	DBMigrationsPath = "DbMigrationPath"
	JWTSignature     = "JwtSignature"
	JWTAccessTTL     = "JwtAccessTtl"
	JWTRefreshTTL    = "JwtRefreshTtl"
	WebSocketAddress = "WebSocketAddress"
	SSLCertPath      = "SSLCertPath"
	SSLKeyPath       = "SSLKeyPath"
	MailLogin        = "MailLogin"
	MailPassword     = "MailPassword"
	MailHost         = "MailHost"
	MailPort         = "MailPort"
	GRPCAddress      = "GrpcAddress"
)

type DB struct {
	DSN            string
	MigrationsPath string
}

type JWT struct {
	Signature  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type WebSocket struct {
	Address string
}

type SSL struct {
	CertPath string
	KeyPath  string
}

type Mail struct {
	Login    string
	Password string
	Host     string
	Port     string
}

type GRPC struct {
	Address string
}

type AppConfig struct {
	DB        *DB
	JWT       *JWT
	WebSocket *WebSocket
	SSL       *SSL
	Mail      *Mail
	GRPC      *GRPC
}

func ReadConfig() *AppConfig {
	setDefaults()
	readEnv()
	readFlags()

	return &AppConfig{
		DB: &DB{
			DSN:            viper.GetString(DBDsn),
			MigrationsPath: viper.GetString(DBMigrationsPath),
		},
		JWT: &JWT{
			Signature:  viper.GetString(JWTSignature),
			AccessTTL:  viper.GetDuration(JWTAccessTTL),
			RefreshTTL: viper.GetDuration(JWTRefreshTTL),
		},
		WebSocket: &WebSocket{
			Address: viper.GetString(WebSocketAddress),
		},
		SSL: &SSL{
			CertPath: viper.GetString(SSLCertPath),
			KeyPath:  viper.GetString(SSLKeyPath),
		},
		Mail: &Mail{
			Login:    viper.GetString(MailLogin),
			Password: viper.GetString(MailPassword),
			Host:     viper.GetString(MailHost),
			Port:     viper.GetString(MailPort),
		},
		GRPC: &GRPC{
			Address: viper.GetString(GRPCAddress),
		},
	}
}

func readEnv() {
	_ = viper.BindEnv(DBDsn, "DB_DSN")
	_ = viper.BindEnv(DBMigrationsPath, "DB_MIGRATIONS_PATH")
	_ = viper.BindEnv(JWTSignature, "JWT_SIGNATURE")
	_ = viper.BindEnv(JWTAccessTTL, "JWT_ACCESS_TTL")
	_ = viper.BindEnv(JWTRefreshTTL, "JWT_REFRESH_TTL")
	_ = viper.BindEnv(WebSocketAddress, "WEBSOCKET_ADDRESS")
	_ = viper.BindEnv(SSLCertPath, "SSL_CERT_PATH")
	_ = viper.BindEnv(SSLKeyPath, "SSL_KEY_PATH")
	_ = viper.BindEnv(MailLogin, "MAIL_LOGIN")
	_ = viper.BindEnv(MailPassword, "MAIL_PASSWORD")
	_ = viper.BindEnv(MailHost, "MAIL_HOST")
	_ = viper.BindEnv(MailPort, "MAIL_PORT")
	_ = viper.BindEnv(GRPCAddress, "GRPC_ADDRESS")
}

func readFlags() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	pflag.String("db-dsn", "", "DB DSN")
	pflag.String("db-mig-path", "", "Path to db migrations")
	pflag.String("jwt-sig", "", "Secret jwt signature")
	pflag.Duration("jwt-attl", 60*time.Second, "TTL for JWT access token")
	pflag.Duration("jwt-rttl", 360*time.Second, "TTL for JWT refresh token")
	pflag.String("ws-addr", "", "Websocket server address")
	pflag.String("ssl-cert", "", "Path to ssl cert")
	pflag.String("ssl-key", "", "Path to ssl key")
	pflag.String("mail-login", "", "Mail login")
	pflag.String("mail-pass", "", "Mail password")
	pflag.String("mail-host", "", "Mail host")
	pflag.String("mail-port", "", "Mail port")
	pflag.String("grpc-addr", "", "GRPC server address")

	pflag.Parse()

	_ = viper.BindPFlag(DBDsn, pflag.Lookup("db-dsn"))
	_ = viper.BindPFlag(DBMigrationsPath, pflag.Lookup("db-mig-path"))
	_ = viper.BindPFlag(JWTSignature, pflag.Lookup("jwt-sig"))
	_ = viper.BindPFlag(JWTAccessTTL, pflag.Lookup("jwt-attl"))
	_ = viper.BindPFlag(JWTRefreshTTL, pflag.Lookup("jwt-rttl"))
	_ = viper.BindPFlag(WebSocketAddress, pflag.Lookup("ws-addr"))
	_ = viper.BindPFlag(SSLCertPath, pflag.Lookup("ssl-cert"))
	_ = viper.BindPFlag(SSLKeyPath, pflag.Lookup("ssl-key"))
	_ = viper.BindPFlag(MailLogin, pflag.Lookup("mail-login"))
	_ = viper.BindPFlag(MailPassword, pflag.Lookup("mail-pass"))
	_ = viper.BindPFlag(MailHost, pflag.Lookup("mail-host"))
	_ = viper.BindPFlag(MailPort, pflag.Lookup("mail-port"))
	_ = viper.BindPFlag(GRPCAddress, pflag.Lookup("grpc-addr"))
}

func setDefaults() {
	viper.SetDefault(DBDsn, "postgres://postgres:postgres@localhost:5432/gophkeeper")
	viper.SetDefault(DBMigrationsPath, "internal/server/infrastructure/migrations")
	viper.SetDefault(JWTSignature, "secret")
	viper.SetDefault(JWTAccessTTL, 60*time.Second)
	viper.SetDefault(JWTRefreshTTL, 360*time.Second)
	viper.SetDefault(WebSocketAddress, ":7777")
	viper.SetDefault(SSLCertPath, "ca.crt")
	viper.SetDefault(SSLKeyPath, "ca.key")
	viper.SetDefault(MailLogin, "from@gmail.com")
	viper.SetDefault(MailPassword, "")
	viper.SetDefault(MailHost, "localhost")
	viper.SetDefault(MailPort, "1025")
	viper.SetDefault(GRPCAddress, ":3200")
}
