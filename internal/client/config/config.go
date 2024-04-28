package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	JWTSignature               = "JwtSignature"
	JWTAccessTTL               = "JwtAccessTtl"
	JWTRefreshTTL              = "JwtRefreshTtl"
	WebSocketAddress           = "WebSocketAddress"
	WebSocketConnectionTimeout = "WebSocketConnectionTimeout"
	GRPCAddress                = "GrpcAddress"
)

type JWT struct {
	Signature  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type WebSocket struct {
	ServerAddress     string
	ConnectionTimeout time.Duration
}

type GRPC struct {
	ServerAddress string
}

type AppConfig struct {
	JWT       *JWT
	WebSocket *WebSocket
	GRPC      *GRPC
}

func ReadConfig() *AppConfig {
	setDefaults()
	readEnv()
	readFlags()

	return &AppConfig{
		JWT: &JWT{
			Signature:  viper.GetString(JWTSignature),
			AccessTTL:  viper.GetDuration(JWTAccessTTL),
			RefreshTTL: viper.GetDuration(JWTRefreshTTL),
		},
		WebSocket: &WebSocket{
			ServerAddress:     viper.GetString(WebSocketAddress),
			ConnectionTimeout: viper.GetDuration(WebSocketConnectionTimeout),
		},
		GRPC: &GRPC{
			ServerAddress: viper.GetString(GRPCAddress),
		},
	}
}

func readEnv() {
	_ = viper.BindEnv(JWTSignature, "JWT_SIGNATURE")
	_ = viper.BindEnv(JWTAccessTTL, "JWT_ACCESS_TTL")
	_ = viper.BindEnv(JWTRefreshTTL, "JWT_REFRESH_TTL")
	_ = viper.BindEnv(WebSocketAddress, "WEBSOCKET_ADDRESS")
	_ = viper.BindEnv(WebSocketConnectionTimeout, "WEBSOCKET_CONNECTION_TIMEOUT")
	_ = viper.BindEnv(GRPCAddress, "GRPC_ADDRESS")
}

func readFlags() {

	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)

	pflag.String("jwt-sig", "", "Secret jwt signature")
	pflag.Duration("jwt-attl", 60*time.Second, "TTL for JWT access token")
	pflag.Duration("jwt-rttl", 360*time.Second, "TTL for JWT refresh token")
	pflag.String("ws-addr", "", "Websocket server address")
	pflag.Duration("ws-ttl", 100*time.Millisecond, "Timeout to connect to websocket server")
	pflag.String("grpc-addr", "", "GRPC server address")

	pflag.Parse()

	_ = viper.BindPFlag(JWTSignature, pflag.Lookup("jwt-sig"))
	_ = viper.BindPFlag(JWTAccessTTL, pflag.Lookup("jwt-attl"))
	_ = viper.BindPFlag(JWTRefreshTTL, pflag.Lookup("jwt-rttl"))
	_ = viper.BindPFlag(WebSocketAddress, pflag.Lookup("ws-addr"))
	_ = viper.BindPFlag(WebSocketConnectionTimeout, pflag.Lookup("ws-ttl"))
	_ = viper.BindPFlag(GRPCAddress, pflag.Lookup("grpc-addr"))
}

func setDefaults() {
	viper.SetDefault(JWTSignature, "secret")
	viper.SetDefault(JWTAccessTTL, 60*time.Second)
	viper.SetDefault(JWTRefreshTTL, 360*time.Second)
	viper.SetDefault(WebSocketAddress, ":7777")
	viper.SetDefault(WebSocketConnectionTimeout, 100*time.Millisecond)
	viper.SetDefault(GRPCAddress, ":3200")
}
