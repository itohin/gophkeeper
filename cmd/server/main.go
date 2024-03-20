package main

import (
	"context"
	"github.com/itohin/gophkeeper/internal/server/adapters/db/postgres"
	"github.com/itohin/gophkeeper/internal/server/adapters/grpc"
	"github.com/itohin/gophkeeper/internal/server/usecases/auth"
	"github.com/itohin/gophkeeper/internal/server/usecases/secrets"
	"github.com/itohin/gophkeeper/pkg/database"
	"github.com/itohin/gophkeeper/pkg/hash/password"
	"github.com/itohin/gophkeeper/pkg/jwt"
	"github.com/itohin/gophkeeper/pkg/logger"
	"github.com/itohin/gophkeeper/pkg/mailer"
	"github.com/itohin/gophkeeper/pkg/otp"
	"github.com/itohin/gophkeeper/pkg/uuid"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	l := logger.NewLogger()
	db, err := database.NewPgxPoolDB(
		context.Background(),
		"postgres://postgres:postgres@localhost:5432/gophkeeper",
		"internal/server/infrastructure/migrations",
	)
	if err != nil {
		l.Fatal(err)
	}
	defer db.Pool.Close()

	srv, err := setupServer(db, l)
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
		close(idleConnsClosed)
	}()

	err = srv.Start()
	if err != nil {
		l.Fatal(err)
	}
	<-idleConnsClosed

	db.Pool.Close()

	l.Info("server shutdown gracefully")

}

func setupServer(db *database.PgxPoolDB, l logger.Logger) (*grpc.Server, error) {
	authUseCase, err := setupAuth(db, l)
	secretsUseCase := secrets.NewSecretsUseCase()
	if err != nil {
		return nil, err
	}
	return grpc.NewServer(authUseCase, secretsUseCase, l), nil
}

func setupAuth(db *database.PgxPoolDB, l logger.Logger) (*auth.AuthUseCase, error) {
	usersRepo := postgres.NewUsersRepository(db)
	sessionsRepo := postgres.NewSessionsRepository(db)
	tx := database.NewPgxTransaction(db.Pool)
	uuidGen := uuid.NewGoogleUUIDGenerator()
	jwtGen, err := jwt.NewJWTGOManager("secret", 60*time.Second, 360*time.Second)
	if err != nil {
		return nil, err
	}
	passwordHash := password.NewBcryptPasswordHasher()
	otpGen := otp.NewGOTPGenerator(9)
	smtp := mailer.NewSMTPMailer("from@gmail.com", "", "localhost", "1025", l)

	return auth.NewAuthUseCase(passwordHash, uuidGen, otpGen, usersRepo, sessionsRepo, smtp, jwtGen, tx), nil
}
