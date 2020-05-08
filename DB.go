package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

var DBConnection *pgx.Conn

//postgres://postgres:1337@/anonchat-tgbot?host=/cloudsql/tg-bots-276110:europe-west6:tgbots-db
func init() {
	//var err error
	//
	//DBConnection, err = pgx.Connect(context.Background(), "postgres://postgres:1337@34.65.65.169:5432/anonchat-tgbot")
	//if err != nil {
	//	DBConnectionError(err)
	//}
	//
	//log.Println("Connected to PSQL!")
}

func UserFirstStart(userId int) {
	x, err := DBConnection.Exec(context.Background(), "INSERT INTO users VALUES($1, $2, $3)",
		userId, false, false)

	if err != nil {
		DBQueryError(err)
	}
	x.Insert()
}

func CheckUserReg(userId int) (isRegistered bool) {
	var findedUser int

	err := DBConnection.QueryRow(context.Background(), "SELECT user_id FROM users WHERE user_id = $1",
		userId).Scan(&findedUser)
	if err != nil {
		DBQueryError(err)
	}

	if findedUser != 0 {
		isRegistered = true
	}

	return
}

func IsUserChatting(userId int) (isChat bool) {
	err := DBConnection.QueryRow(context.Background(), "SELECT is_chatting FROM users WHERE user_id = $1",
		userId).Scan(&isChat)
	if err != nil {
		DBQueryError(err)
	}

	return
}

func ChangeUserChattingState(userId int, status bool) {
	x, err := DBConnection.Exec(context.Background(), "UPDATE users SET is_chatting = $2 WHERE user_id = $1",
		userId, status)
	if err != nil {
		DBQueryError(err)
	}
	x.Update()
}

func IsUserSearching(userId int) (isSearch bool) {
	err := DBConnection.QueryRow(context.Background(), "SELECT is_searching FROM users WHERE user_id = $1",
		userId).Scan(&isSearch)
	if err != nil {
		DBQueryError(err)
	}

	return
}

func ChangeUserSearchingState(userId int, status bool) {
	x, err := DBConnection.Exec(context.Background(), "UPDATE users SET is_searching = $2 WHERE user_id = $1",
		userId, status)
	if err != nil {
		DBQueryError(err)
	}
	x.Update()
}

func FindFreeUsers() (freeUsers []int) {
	activeUsers, err := DBConnection.Query(context.Background(),
		"SELECT user_id FROM users WHERE is_chatting = false AND is_searching = true")
	if err != nil {
		DBQueryError(err)
	}

	defer activeUsers.Close()

	var oneUser int
	for activeUsers.Next() {
		err := activeUsers.Scan(&oneUser)
		if err != nil {
			DBScanError(err)
		}
		freeUsers = append(freeUsers, oneUser)
	}

	return
}

func AddNewChat(firstUserId int, secondUserId int) {
	x, err := DBConnection.Exec(context.Background(), "INSERT INTO chats VALUES($1, $2)",
		firstUserId, secondUserId)
	if err != nil {
		DBQueryError(err)
	}
	x.Insert()
}

func FindSecondUserFromChat(userId int) (secondUser int) {
	err := DBConnection.QueryRow(context.Background(),
		"SELECT second_user FROM chats WHERE first_user = $1", userId).Scan(&secondUser)
	if err != nil {
		DBQueryError(err)
	}

	return
}

//func DeleteChat(userId int) {
//	x, err := DBConnection.Exec(context.Background(), "DELETE FROM chats WHERE first_user = $1",
//		userId)
//	if err != nil {
//		DBQueryError(err)
//	}
//	x.Delete()
//}

func DBScanError(err error) {
	log.Println(fmt.Errorf("Scan failed: %w\n", err))
}

func DBQueryError(err error) {
	log.Println(fmt.Errorf("QueryRow failed: %v \n", err))
}

func DBConnectionError(err error) {
	log.Println(fmt.Errorf("Unable to connection to database: %w\n", err))
	panic(err)
}
