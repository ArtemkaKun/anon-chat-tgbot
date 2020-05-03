package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var DBConnection *pgx.Conn

func init() {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1337@34.65.183.194:5432/anonchat-tgbot")
	if err != nil {
		DBConnectionError(err)
	}

	fmt.Println("Connected to PSQL!")

	DBConnection = connection
}

func UserFirstStart(user_id int) {
	_, err := DBConnection.Exec(context.Background(), "INSERT INTO users VALUES($1, $2, $3)",
		user_id, false, false)
	if err != nil {
		DBQueryError(err)
	}
}

func CheckUserReg(user_id int) bool {
	regStatus, err := DBConnection.Query(context.Background(), "SELECT user_id FROM users WHERE user_id = $1",
		user_id)
	if err != nil {
		DBQueryError(err)
	}

	defer regStatus.Close()

	return regStatus.Next()
}

func IsUserChatting(user_id int) bool {
	chat_status, err := DBConnection.Query(context.Background(), "SELECT is_chatting FROM users WHERE user_id = $1",
		user_id)
	if err != nil {
		DBQueryError(err)
	}

	defer chat_status.Close()

	var is_chat bool

	for chat_status.Next() {
		err := chat_status.Scan(&is_chat)
		if err != nil {
			DBScanError(err)
		}
	}

	if chat_status.Err() != nil {
		DBQueryError(err)
	}

	return is_chat
}

func ChangeUserChattingState(user_id int, status bool) {
	_, err := DBConnection.Exec(context.Background(), "UPDATE users SET is_chatting = $2 WHERE user_id = $1",
		user_id, status)
	if err != nil {
		DBQueryError(err)
	}
}

func IsUserSearching(user_id int) bool {
	search_status, err := DBConnection.Query(context.Background(), "SELECT is_searching FROM users WHERE user_id = $1",
		user_id)
	if err != nil {
		DBQueryError(err)
	}

	defer search_status.Close()

	var is_search bool

	for search_status.Next() {
		err := search_status.Scan(&is_search)
		if err != nil {
			DBScanError(err)
		}
	}

	if search_status.Err() != nil {
		DBQueryError(err)
	}

	return is_search
}

func ChangeUserSearchingState(user_id int, status bool) {
	_, err := DBConnection.Exec(context.Background(), "UPDATE users SET is_searching = $2 WHERE user_id = $1",
		user_id, status)
	if err != nil {
		DBQueryError(err)
	}
}

func FindFreeUsers() []int {
	active_users, err := DBConnection.Query(context.Background(),
		"SELECT user_id FROM users WHERE is_chatting = false AND is_searching = true")
	if err != nil {
		DBQueryError(err)
	}

	defer active_users.Close()

	var free_users []int
	var one_user int

	for active_users.Next() {
		err := active_users.Scan(&one_user)
		if err != nil {
			DBScanError(err)
		}
		free_users = append(free_users, one_user)
	}

	return free_users
}

func AddNewChat(first_user_id int, second_user_id int) {
	_, err := DBConnection.Exec(context.Background(), "INSERT INTO chats VALUES($1, $2)",
		first_user_id, second_user_id)
	if err != nil {
		DBQueryError(err)
	}
}

func FindSecondUserFromChat(user_id int) int {
	next_chat_user, err := DBConnection.Query(context.Background(),
		"SELECT second_user FROM chats WHERE first_user = $1", user_id)
	if err != nil {
		DBQueryError(err)
	}

	defer next_chat_user.Close()

	var second_user = 0

	for next_chat_user.Next() {
		err := next_chat_user.Scan(&second_user)
		if err != nil {
			DBScanError(err)
		}
	}

	if next_chat_user.Err() != nil {
		DBQueryError(err)
	}

	return second_user
}

func DeleteChat(user_id int) {
	_, err := DBConnection.Exec(context.Background(), "DELETE FROM chats WHERE first_user = $1",
		user_id)
	if err != nil {
		DBQueryError(err)
	}
}

func DBScanError(err error) {
	fmt.Println(fmt.Errorf("Scan failed: %w\n", err))
}

func DBQueryError(err error) {
	fmt.Println(fmt.Errorf("QueryRow failed: %v \n", err))
}

func DBConnectionError(err error) {
	fmt.Println(fmt.Errorf("Unable to connection to database: %w\n", err))
	panic(err)
}
