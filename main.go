package main

import (
	"database/sql"
	"github.com/Syfaro/telegram-bot-api"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func main() {
	bot := BotStart()
	my_db := DBStart()
	defer my_db.Close()
	BotUpdateLoop(bot, my_db)
}

func BotStart() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI("1057128816:AAE3MrZxSXnMPV1UNYuLbOQobd-sxUIhGw4")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Autorised on account %s", bot.Self.UserName)

	return bot
}
func BotUpdateLoop(my_bot *tgbotapi.BotAPI, database *sql.DB) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := my_bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			if FindChat(update.Message.From.ID, database) != 0{
				chat_id := FindChat(update.Message.From.ID, database)
				msg := tgbotapi.NewMessage(int64(chat_id), "Bot cannot send photos, stickers, documents, audios and videos yet!")
				//photo := tgbotapi.NewPhotoUpload(int64(chat_id), nil)

				//photo.File = update.Message.Photo
				//photo.UseExisting = true
				msg.Text = update.Message.Text
				if msg.Text != "" {
					_, err := my_bot.Send(msg)
					if err != nil {
						log.Panic(err)
					}
				} else {
					_, err := my_bot.Send(msg)
					if err != nil {
						log.Panic(err)
					}
				}


				continue
			} else {
				continue
			}
		}

		chat_id := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chat_id, "")

		switch update.Message.Command() {
		case "start":
			if !CheckReg(update.Message.From.ID, database) {
				FirstStart(update.Message.From.ID, database)
				msg.Text = "You was registered"
			} else {
				msg.Text = "You are already registered"
			}
		case "go_chat":
			if CheckReg(update.Message.From.ID, database) {
				if IsFree(update.Message.From.ID, database) {
					if !IsSearch(update.Message.From.ID, database) {
						go SearchChat(update, database, update.Message.From.ID, int64(update.Message.From.ID), my_bot)
						msg.Text = "Search started"
					} else {
						msg.Text = "You are already search"
					}
				} else {
					msg.Text = "You are already chatting"
				}
			} else {
				msg.Text = "You need to register first"
			}
		}

		_, err := my_bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	}
}

func SearchChat(update tgbotapi.Update, database *sql.DB, user_id int, chat_id int64, my_bot *tgbotapi.BotAPI) {
	free_user := 0
	msg := tgbotapi.NewMessage(chat_id, "")

	ChangeSearch(database, user_id, 1)

	for free_user == 0 {
		free_user = FindFree(update.Message.From.ID, database)
		amt := time.Duration(1000)
		time.Sleep(time.Millisecond * amt)
		continue
	}

	ChangeState(database, update.Message.From.ID, 1)
	ChangeSearch(database, user_id, 0)
	AddChat(update.Message.From.ID, free_user, database)

	msg.Text = "Now you can chat"

	_, err := my_bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}
func DBStart() *sql.DB {
	my_db, err := sql.Open("mysql", "root:1337@/anonstudchat")
	if err != nil {
		log.Panic(err)
	} else {
		err = my_db.Ping()
		if err != nil {
			log.Panic(err)
		}
	}
	return my_db
}
func FirstStart(user_id int, my_db *sql.DB) {
	stmtIns, err := my_db.Prepare("INSERT INTO users_info VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(user_id, 0, 0)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}
func CheckReg(user_id int, my_db *sql.DB) bool {
	stmtOut, err := my_db.Prepare("SELECT user_id FROM users_info WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}

	var is_reg int
	err = stmtOut.QueryRow(user_id).Scan(&is_reg)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		panic(err.Error())
	}

	if is_reg != 0 {
		return true
	} else {
		return false
	}

}
func IsFree(user_id int, my_db *sql.DB) bool {
	stmtOut, err := my_db.Prepare("SELECT user_free FROM users_info WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}

	var is_free int
	err = stmtOut.QueryRow(user_id).Scan(&is_free)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		panic(err.Error())
	}

	if is_free == 0 {
		return true
	} else {
		return false
	}

}
func IsSearch(user_id int, my_db *sql.DB) bool {
	stmtOut, err := my_db.Prepare("SELECT is_search FROM users_info WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}

	var is_free int
	err = stmtOut.QueryRow(user_id).Scan(&is_free)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			panic(err.Error())
		}
		return false
	}

	err = stmtOut.Close()
	if err != nil {
		panic(err.Error())
	}

	if is_free == 1 {
		return true
	} else {
		return false
	}

}
func FindFree(user_id int, my_db *sql.DB) int {
	stmtOut, err := my_db.Prepare("SELECT user_id FROM users_info WHERE user_free = 0 AND is_search = 1 AND user_id != ?")
	if err != nil {
		panic(err.Error())
	}

	var user_free int
	err = stmtOut.QueryRow(user_id).Scan(&user_free)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			panic(err.Error())
		}
		return 0
	}

	err = stmtOut.Close()
	if err != nil {
		panic(err.Error())
	}

	return user_free
}
func ChangeSearch(my_db *sql.DB, user_id int, status int) {
	stmtIns, err := my_db.Prepare("UPDATE users_info SET is_search = ? WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(status, user_id)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}
func ChangeState(my_db *sql.DB, user_id int, status int) {
	stmtIns, err := my_db.Prepare("UPDATE users_info SET user_free = ? WHERE user_id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(status, user_id)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}
func AddChat(first_user_id int, second_user_id int, my_db *sql.DB) {
	stmtIns, err := my_db.Prepare("INSERT INTO chat_buffer VALUES (?, ?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(first_user_id, second_user_id)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}
func FindChat(user_id int, my_db *sql.DB) int {
	stmtOut, err := my_db.Prepare("SELECT second_user FROM chat_buffer WHERE first_user = ?")
	if err != nil {
		panic(err.Error())
	}

	var second_user int
	err = stmtOut.QueryRow(user_id).Scan(&second_user)
	if err != nil {
		err = stmtOut.Close()
		if err != nil {
			panic(err.Error())
		}
		return 0
	}

	err = stmtOut.Close()
	if err != nil {
		panic(err.Error())
	}

	return second_user
}
