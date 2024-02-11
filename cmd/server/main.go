package main

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/server/adapters/db/postgres"
	"github.com/itohin/gophkeeper/pkg/database"
	"github.com/itohin/gophkeeper/pkg/logger"
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

	usersRepo := postgres.NewUsersRepository(db)
	user, err := usersRepo.FindByID(context.Background(), "1955a7d6-0968-425b-bdb6-fb9a0e4b39e7")
	if err != nil {
		fmt.Printf("save error: %v", err)
	}
	fmt.Printf("user: %v", user)
}
