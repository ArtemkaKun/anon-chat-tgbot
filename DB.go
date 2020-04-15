package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var DBConnection *pgx.Conn

func init() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/anonchat-tgbot")
	if err != nil {
		DBConnectionError(err)
		return
	}

	fmt.Println("Connected to PSQL!")

	DBConnection = connection
}

func DBConnectionError(err error) {
	fmt.Println(fmt.Errorf("Unable to connection to database: %w\n", err))
}
