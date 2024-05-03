package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itohin/gophkeeper/internal/server/adapters/db/hydrator"
	"github.com/itohin/gophkeeper/internal/server/adapters/db/postgres"
	"github.com/itohin/gophkeeper/internal/server/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/server/adapters/websocket"
	"github.com/itohin/gophkeeper/internal/server/config"
	"github.com/itohin/gophkeeper/internal/server/usecases/auth"
	"github.com/itohin/gophkeeper/internal/server/usecases/secrets"
	"github.com/itohin/gophkeeper/pkg/database"
	"github.com/itohin/gophkeeper/pkg/events"
	"github.com/itohin/gophkeeper/pkg/hash/password"
	"github.com/itohin/gophkeeper/pkg/jwt"
	"github.com/itohin/gophkeeper/pkg/logger"
	"github.com/itohin/gophkeeper/pkg/mailer"
	"github.com/itohin/gophkeeper/pkg/otp"
	"github.com/itohin/gophkeeper/pkg/uuid"
)

func main() {

	l := logger.NewLogger()

	cfg := config.ReadConfig()

	db, err := database.NewPgxPoolDB(context.Background(), cfg.DB.DSN, cfg.DB.MigrationsPath)
	if err != nil {
		l.Fatal(err)
	}
	defer db.Pool.Close()

	jwtManager, err := jwt.NewJWTGOManager(cfg.JWT.Signature, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	if err != nil {
		l.Fatal(err)
	}

	secretEventsCh := make(chan *events.SecretEvent, 10)
	ws := websocket.NewWSNotifier(cfg.WebSocket.Address, cfg.SSL.CertPath, cfg.SSL.KeyPath, secretEventsCh)

	srv, err := setupServer(db, l, jwtManager, secretEventsCh, cfg)
	if err != nil {
		l.Fatal(err)
	}

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Stop(ctx)
		ws.Stop(ctx)
		close(idleConnsClosed)
	}()
	go ws.Run()

	err = srv.Start()
	if err != nil {
		l.Fatal(err)
	}

	<-idleConnsClosed

	db.Pool.Close()

	l.Info("server shutdown gracefully")

}

func setupServer(db *database.PgxPoolDB, l logger.Logger, jm *jwt.JWTGOManager, eventCh chan *events.SecretEvent, cfg *config.AppConfig) (*grpc.Server, error) {
	uuidGen := uuid.NewGoogleUUIDGenerator()
	secretsRepo := postgres.NewSecretsRepository(db)

	authUseCase, err := setupAuth(db, l, jm, uuidGen, cfg)
	secretsUseCase := secrets.NewSecretsUseCase(uuidGen, secretsRepo, eventCh)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(authUseCase, secretsUseCase, l, jm, hydrator.NewSecretsHydrator(), cfg), nil
}

func setupAuth(db *database.PgxPoolDB, l logger.Logger, jm *jwt.JWTGOManager, uuidGen *uuid.GoogleUUIDGenerator, cfg *config.AppConfig) (*auth.AuthUseCase, error) {
	usersRepo := postgres.NewUsersRepository(db)
	sessionsRepo := postgres.NewSessionsRepository(db)
	tx := database.NewPgxTransaction(db.Pool)
	passwordHash := password.NewBcryptPasswordHasher()
	otpGen := otp.NewGOTPGenerator(9)
	smtp := mailer.NewSMTPMailer(cfg.Mail.Login, cfg.Mail.Password, cfg.Mail.Host, cfg.Mail.Port, l)

	return auth.NewAuthUseCase(passwordHash, uuidGen, otpGen, usersRepo, sessionsRepo, smtp, jm, tx), nil
}
