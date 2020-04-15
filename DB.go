package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var DBConnection *pgx.Conn

func init() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1337@localhost:5432/anonchat-tgbot")
	if err != nil {
		DBConnectionError(err)
	}

	fmt.Println("Connected to PSQL!")

	DBConnection = connection
}

func UserFirstStart(user_id int) {
	_, err := DBConnection.Exec(context.Background(), "INSERT INTO public.users VALUES($1, $2, $3)",
		user_id, false, false)
	if err != nil {
		DBQueryError(err)
	}
}

func CheckUserReg(user_id int) bool {
	regStatus, err := DBConnection.Query(context.Background(), "SELECT user_id FROM public.users WHERE user_id = $1",
		user_id)
	if err != nil {
		DBQueryError(err)
	}

	defer regStatus.Close()

	return regStatus.Next()
}

func DBQueryError(err error) {
	fmt.Println(fmt.Errorf("QueryRow failed: %w\n", err))
}

func DBConnectionError(err error) {
	fmt.Println(fmt.Errorf("Unable to connection to database: %w\n", err))
	panic(err)
}
