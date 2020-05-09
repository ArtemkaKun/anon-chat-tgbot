package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"math/rand"
	"time"
)

const (
	PAUSE_FOR_BACKUPS          = 100
	PAUSE_FOR_CHATMAKER        = 1
	USERS_AMOUNT_FOR_CHATMAKER = 2
)

func main() {
	go BackupCache()
	go ChatMaker()
	BotUpdateLoop()
}

func BackupCache() {
	for {
		Pause(PAUSE_FOR_BACKUPS)

		log.Println("Backup")
		BackupData(Users, Chats, Rooms)
	}
}

func ChatMaker() {
	for true {
		freeUsers := SearchingUsersList()
		usersAmount := len(freeUsers)

		if usersAmount >= USERS_AMOUNT_FOR_CHATMAKER {
			CreateNewChat(usersAmount, freeUsers)
		}

		Pause(PAUSE_FOR_CHATMAKER)
	}
}

func CreateNewChat(usersAmount int, freeUsers []int64) {
	firstUser, secondUser := ChooseRandomUsers(usersAmount, freeUsers)

	MakeChat(firstUser, secondUser)

	for _, user := range []int64{firstUser, secondUser} {
		BotSendMessage(tgbotapi.NewMessage(user, "Now you can chat"))
	}
}

func ChooseRandomUsers(usersAmount int, freeUsers []int64) (firstUser, secondUser int64) {
	rand.Seed(time.Now().UnixNano())
	userOne := rand.Intn(usersAmount)
	userTwo := rand.Intn(usersAmount)

	for userTwo == userOne {
		userTwo = rand.Intn(usersAmount)
	}

	firstUser = freeUsers[userOne]
	secondUser = freeUsers[userTwo]
	return
}

func MakeChat(firstUser int64, secondUser int64) {
	ChangeUserSearchingStatus(false, firstUser, secondUser)
	ChangeUserChattingStatus(true, firstUser, secondUser)

	AddChat(firstUser, secondUser)
}

func Pause(seconds uint16) {
	amt := time.Duration(seconds)
	time.Sleep(time.Second * amt)
}
