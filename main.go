package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"math/rand"
	"time"
)

func main() {
	//go ChatMaker()
	BotUpdateLoop()
}

func ChatMaker() {
	for true {
		freeUsers := SearchingUsersList()
		usersAmount := len(freeUsers)

		if usersAmount > 1 {
			rand.Seed(time.Now().UnixNano())
			userOne := rand.Intn(usersAmount)
			userTwo := rand.Intn(usersAmount)

			for userTwo == userOne {
				userTwo = rand.Intn(usersAmount)
			}

			firstUser := freeUsers[userOne]
			secondUser := freeUsers[userTwo]

			MakeChat(firstUser, secondUser)

			msg := tgbotapi.NewMessage(firstUser, "Now you can chat")
			BotSendMessage(msg)

			msg = tgbotapi.NewMessage(secondUser, "Now you can chat")
			BotSendMessage(msg)
		}

		amt := time.Duration(1000)
		time.Sleep(time.Millisecond * amt)
	}
}

func MakeChat(firstUser int64, secondUser int64) {
	ChangeUserSearchingStatus(firstUser, false)
	ChangeUserSearchingStatus(secondUser, false)

	ChangeUserChattingStatus(firstUser, true)
	ChangeUserChattingStatus(secondUser, true)

	AddChat(firstUser, secondUser)
}
