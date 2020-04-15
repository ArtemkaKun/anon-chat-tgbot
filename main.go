package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"math/rand"
	"time"
)

func main() {
	go ChatMaker()
	BotUpdateLoop()
}

func ChatMaker() {
	for true {
		free_users := FindFreeUsers()
		users_amount := len(free_users)

		if users_amount > 1 {
			rand.Seed(time.Now().UnixNano())
			user_one := rand.Intn(users_amount)
			user_two := rand.Intn(users_amount)

			for user_two == user_one {
				user_two = rand.Intn(users_amount)
			}

			first_user := free_users[user_one]
			second_user := free_users[user_two]

			MakeChat(first_user, second_user)

			msg := tgbotapi.NewMessage(int64(first_user), "Now you can chat")
			BotSendMessage(msg)

			msg = tgbotapi.NewMessage(int64(second_user), "Now you can chat")
			BotSendMessage(msg)
		}

		amt := time.Duration(1000)
		time.Sleep(time.Millisecond * amt)
	}
}

func MakeChat(first_user int, second_user int) {
	ChangeUserSearchingState(first_user, false)
	ChangeUserSearchingState(second_user, false)

	ChangeUserChattingState(first_user, true)
	ChangeUserChattingState(second_user, true)

	AddNewChat(first_user, second_user)
	AddNewChat(second_user, first_user)
}
