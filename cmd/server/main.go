package main

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/server/adapters/db/postgres"
	auth2 "github.com/itohin/gophkeeper/internal/server/usecases/auth"
	"github.com/itohin/gophkeeper/pkg/database"
	"github.com/itohin/gophkeeper/pkg/hash/password"
	"github.com/itohin/gophkeeper/pkg/jwt"
	"github.com/itohin/gophkeeper/pkg/logger"
	mailer2 "github.com/itohin/gophkeeper/pkg/mailer"
	"github.com/itohin/gophkeeper/pkg/otp"
	"github.com/itohin/gophkeeper/pkg/uuid"
	"time"
)

type PGXTransaction struct {
}

func NewPGXTransaction() *PGXTransaction {
	return &PGXTransaction{}
}

func (t *PGXTransaction) Transaction(f func() error) error {
	defer fmt.Println("defer call")
	return f()
}

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

	usersRepo := postgres.NewUsersRepository(db)
	sessionsRepo := postgres.NewSessionsRepository(db)
	tx := database.NewPgxTransaction(db.Pool)
	uuidGen := uuid.NewGoogleUUIDGenerator()
	jwtGen, err := jwt.NewJWTGOManager("secret", 60*time.Second, 360*time.Second)
	if err != nil {
		l.Fatal(err)
	}
	passwordHash := password.NewBcryptPasswordHasher()
	otpGen := otp.NewGOTPGenerator(9)
	mailer := mailer2.NewSMTPMailer("from", "pass", "host", "port", l)

	auth := auth2.NewAuthUseCase(passwordHash, uuidGen, otpGen, usersRepo, sessionsRepo, mailer, jwtGen, tx)
	//err = auth.Register(context.Background(), "a@a.com", "password")
	token, err := auth.Verify(context.Background(), "a@a.com", "PLYOGRTHJWCRGMQ", "unique_fingerprint")

	fmt.Println("token: ", token, "error: ", err)

}
