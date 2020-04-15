package main

func main() {
	//go ChatMaker(my_db, bot)
	//BotUpdateLoop(bot, my_db)
}

//func ChatMaker(database *sql.DB, my_bot *tgbotapi.BotAPI) {
//	for true {
//		free_users := FindFree(database, my_bot)
//		users_amount := len(free_users)
//		if users_amount > 1 {
//			rand.Seed(time.Now().UnixNano())
//			first_user := rand.Intn(users_amount)
//			second_user := rand.Intn(users_amount)
//
//			for second_user == first_user {
//				second_user = rand.Intn(users_amount)
//			}
//
//			ChangeSearch(database, free_users[first_user], 0, my_bot)
//			ChangeSearch(database, free_users[second_user], 0, my_bot)
//			ChangeState(database, free_users[first_user], 1, my_bot)
//			ChangeState(database, free_users[second_user], 1, my_bot)
//			AddChat(free_users[first_user], free_users[second_user], database, my_bot)
//			AddChat(free_users[second_user], free_users[first_user], database, my_bot)
//
//			msg := tgbotapi.NewMessage(int64(free_users[first_user]), "")
//			msg.Text = "Now you can chat"
//			_, err := my_bot.Send(msg)
//			if err != nil {
//				ErrorCatch(err.Error(), my_bot)
//			}
//
//			msg = tgbotapi.NewMessage(int64(free_users[second_user]), "")
//			msg.Text = "Now you can chat"
//			_, err = my_bot.Send(msg)
//			if err != nil {
//				ErrorCatch(err.Error(), my_bot)
//			}
//		}
//		amt := time.Duration(1000)
//		time.Sleep(time.Millisecond * amt)
//	}
//}
//
//
//}
//
//func ErrorCatch(err string, my_bot *tgbotapi.BotAPI) {
//	msg := tgbotapi.NewMessage(ADMIN, err)
//	my_bot.Send(msg)
//}
