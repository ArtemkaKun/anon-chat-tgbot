package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var DBConnection = func() (connection *pgx.Conn) {
	var err error

	connection, err = pgx.Connect(context.Background(), os.Getenv("DB_TOKEN"))
	if err != nil {
		DBConnectionError(err)
	}

	log.Println("Connected to PSQL!")

	return
}()

func BackupData(users map[int64]*UserStatuses, chats map[int64]int64, rooms map[string]int64) {
	InsertUsersCache(users)
	InsertChatsCache(chats)
	InsertRoomsCache(rooms)
}

func InsertUsersCache(users map[int64]*UserStatuses) {
	for key, val := range users {
		_, err := DBConnection.Exec(context.Background(), "INSERT INTO users VALUES($1, $2, $3)",
			key, val.IsUserSearching(), val.IsUserChatting())

		if err != nil {
			BackupCacheError(users)
		}
	}
}

func InsertChatsCache(chats map[int64]int64) {
	for key, val := range chats {
		_, err := DBConnection.Exec(context.Background(), "INSERT INTO chats VALUES($1, $2)",
			key, val)

		if err != nil {
			BackupCacheError(chats)
		}
	}
}

func InsertRoomsCache(rooms map[string]int64) {
	for key, val := range rooms {
		_, err := DBConnection.Exec(context.Background(), "INSERT INTO rooms VALUES($1, $2)",
			key, val)

		if err != nil {
			DBQueryError(err)
		}
	}
}

func BackupCacheError(cache interface{}) {
	switch v := cache.(type) {
	case map[int64]*UserStatuses:
		for key, val := range v {
			_, err := DBConnection.Exec(context.Background(), "UPDATE users SET is_searching = $1, is_chatting = $2 WHERE user_id = $3",
				val.IsUserSearching(), val.IsUserChatting(), key)

			if err != nil {
				DBQueryError(err)
			}
		}
	case map[int64]int64:
		for key, val := range v {
			_, err := DBConnection.Exec(context.Background(), "UPDATE chats SET second_user = $1 WHERE first_user = $2",
				val, key)

			if err != nil {
				DBQueryError(err)
			}
		}
	}
}

func GetUsersFromDB() (users map[int64]*UserStatuses) {
	users = make(map[int64]*UserStatuses)

	rows, err := DBConnection.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		DBQueryError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var n int64
		var search, chat bool

		err = rows.Scan(&n, &search, &chat)
		if err != nil {
			DBQueryError(err)
		}

		users[n] = new(UserStatuses)
		users[n].SetSearchingStatus(search)
		users[n].SetChattingStatus(chat)
	}

	return
}

func GetChatsFromDB() (chats map[int64]int64) {
	chats = make(map[int64]int64)

	rows, err := DBConnection.Query(context.Background(), "SELECT * FROM chats")
	if err != nil {
		DBQueryError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var firstUser, secondUser int64

		err = rows.Scan(&firstUser, &secondUser)
		if err != nil {
			DBQueryError(err)
		}

		chats[firstUser] = secondUser
	}

	return
}

func GetRoomsFromDB() (rooms map[string]int64) {
	rooms = make(map[string]int64)

	rows, err := DBConnection.Query(context.Background(), "SELECT * FROM rooms")
	if err != nil {
		DBQueryError(err)
	}

	defer rows.Close()

	for rows.Next() {
		var authorUser int64
		var token string

		err = rows.Scan(&token, &authorUser)
		if err != nil {
			DBQueryError(err)
		}

		rooms[token] = authorUser
	}

	return
}

func DBQueryError(err error) {
	log.Println(fmt.Errorf("QueryRow failed: %v \n", err))
}

func DBConnectionError(err error) {
	log.Println(fmt.Errorf("Unable to connection to database: %w\n", err))
	panic(err)
}
