package main

import (
	"context"
	"github.com/itohin/gophkeeper/internal/server/adapters/db/postgres"
	"github.com/itohin/gophkeeper/pkg/database"
	"log"
)

func main() {
	db, err := database.NewPgxPoolDB(context.Background(), "postgres://postgres:postgres@localhost:5432/gophkeeper")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()

	usersRepo := postgres.NewUsersRepository(db)

	log.Println(usersRepo)
}
