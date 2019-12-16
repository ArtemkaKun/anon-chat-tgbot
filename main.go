package main

import (
	"database/sql"
	"github.com/Syfaro/telegram-bot-api"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"time"
)

func main() {
	bot := BotStart()
	my_db := DBStart()
	defer my_db.Close()
	go ChatMaker(my_db, bot)
	BotUpdateLoop(bot, my_db)
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat"),
	),
)

func BotStart() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI("1022122500:AAFy8sDJFUlgw0e7JURelghBPv_is5kG7ck") //1057128816:AAE3MrZxSXnMPV1UNYuLbOQobd-sxUIhGw4 - AnonStud 1022122500:AAFy8sDJFUlgw0e7JURelghBPv_is5kG7ck - Freedom
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
			switch update.Message.Text {
			case "New chat":
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

				if CheckReg(update.Message.From.ID, database) {
					if IsFree(update.Message.From.ID, database) {
						if !IsSearch(update.Message.From.ID, database) {
							ChangeSearch(database, update.Message.From.ID, 1)
							msg.Text = "Search started"
						} else {
							msg.Text = "You are searching already"
						}
					} else {
						msg.Text = "You are chatting already"
					}
				} else {
					msg.Text = "You need to register first"
				}

				my_bot.Send(msg)
				continue

			case "Leave chat":
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

				if CheckReg(update.Message.From.ID, database) {
					if FindChat(update.Message.From.ID, database) != 0 {
						chat_id := FindChat(update.Message.From.ID, database)
						DeleteChat(update.Message.From.ID, database)
						ChangeState(database, update.Message.From.ID, 0)
						msg.Text = "You leaved a chat"

						DeleteChat(chat_id, database)
						ChangeState(database, chat_id, 0)
						my_bot.Send(tgbotapi.NewMessage(int64(chat_id), "The stranger leave the chat"))
					} else {
						msg.Text = "You are not chatting now!"
					}
				} else {
					msg.Text = "You need to register first"
				}

				my_bot.Send(msg)
				continue
			}

			if FindChat(update.Message.From.ID, database) != 0 {
				chat_id := FindChat(update.Message.From.ID, database)
				msg := tgbotapi.NewMessage(int64(chat_id), "")

				msg.Text = update.Message.Text
				if msg.Text != "" {
					my_bot.Send(msg)
				} else if update.Message.Photo != nil {
					photo := tgbotapi.NewPhotoShare(int64(chat_id), "")
					photo.FileID = (*update.Message.Photo)[2].FileID
					my_bot.Send(photo)
				} else if update.Message.Voice != nil {
					voice := tgbotapi.NewVoiceShare(int64(chat_id), "")
					voice.FileID = update.Message.Voice.FileID
					my_bot.Send(voice)
				} else if update.Message.Animation != nil {
					voice := tgbotapi.NewAnimationShare(int64(chat_id), "")
					voice.FileID = update.Message.Animation.FileID
					my_bot.Send(voice)
				} else if update.Message.Audio != nil {
					audio := tgbotapi.NewAudioShare(int64(chat_id), "")
					audio.FileID = update.Message.Audio.FileID
					my_bot.Send(audio)
				} else if update.Message.Sticker != nil {
					sticker := tgbotapi.NewStickerShare(int64(chat_id), "")
					sticker.FileID = update.Message.Sticker.FileID
					my_bot.Send(sticker)
				} else if update.Message.Document != nil {
					doc := tgbotapi.NewDocumentShare(int64(chat_id), "")
					doc.FileID = update.Message.Document.FileID
					my_bot.Send(doc)
				} else if update.Message.Video != nil {
					video := tgbotapi.NewVideoShare(int64(chat_id), "")
					video.FileID = update.Message.Video.FileID
					my_bot.Send(video)
				} else if update.Message.VideoNote != nil {
					video_note := tgbotapi.NewVideoNoteShare(int64(chat_id),0, "")
					video_note.FileID = update.Message.VideoNote.FileID
					video_note.Length = update.Message.VideoNote.Length
					my_bot.Send(video_note)
				} else {
					msg.Text = "Bot cannot send this shit yet! Please, contact with creator"
					my_bot.Send(msg)
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
				msg.Text = "You were registered"
				msg.ReplyMarkup = numericKeyboard
			} else {
				msg.Text = "You have registered already"
				msg.ReplyMarkup = numericKeyboard
			}
		case "go_chat":
			if CheckReg(update.Message.From.ID, database) {
				if IsFree(update.Message.From.ID, database) {
					if !IsSearch(update.Message.From.ID, database) {
						ChangeSearch(database, update.Message.From.ID, 1)
						msg.Text = "Search started"
					} else {
						msg.Text = "You are searching already"
					}
				} else {
					msg.Text = "You are chatting already"
				}
			} else {
				msg.Text = "You need to register first"
			}
		case "leave_chat":
			if CheckReg(update.Message.From.ID, database) {
				if FindChat(update.Message.From.ID, database) != 0 {
					chat_id := FindChat(update.Message.From.ID, database)
					DeleteChat(update.Message.From.ID, database)
					ChangeState(database, update.Message.From.ID, 0)
					msg.Text = "You leaved a chat"

					DeleteChat(chat_id, database)
					ChangeState(database, chat_id, 0)
					my_bot.Send(tgbotapi.NewMessage(int64(chat_id), "The stranger leave the chat"))
				} else {
					msg.Text = "You are not chatting now!"
				}
			} else {
				msg.Text = "You need to register first"
			}
		}

		my_bot.Send(msg)
	}
}

func ChatMaker( database *sql.DB, my_bot *tgbotapi.BotAPI) {
	for true {
		free_users := FindFree(database)
		users_amount := len(free_users)
		if users_amount > 1 {
			rand.Seed(time.Now().UnixNano())
			first_user := rand.Intn(users_amount)
			second_user := rand.Intn(users_amount)

			for second_user == first_user {
				second_user = rand.Intn(users_amount)
			}

			ChangeSearch(database, free_users[first_user], 0)
			ChangeSearch(database, free_users[second_user], 0)
			ChangeState(database, free_users[first_user], 1)
			ChangeState(database, free_users[second_user], 1)
			AddChat(free_users[first_user], free_users[second_user], database)
			AddChat(free_users[second_user], free_users[first_user], database)

			msg := tgbotapi.NewMessage(int64(free_users[first_user]), "")
			msg.Text = "Now you can chat"
			my_bot.Send(msg)

			msg = tgbotapi.NewMessage(int64(free_users[second_user]), "")
			msg.Text = "Now you can chat"
			my_bot.Send(msg)
		}
		amt := time.Duration(1000)
		time.Sleep(time.Millisecond * amt)
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
func FindFree(my_db *sql.DB) []int {
	stmtOut, err := my_db.Query("SELECT user_id FROM users_info WHERE user_free = 0 AND is_search = 1")
	if err != nil {
		panic(err.Error())
	}

	var user_free []int
	var one_user int

	for stmtOut.Next() {
		err = stmtOut.Scan(&one_user)
		if err != nil {
			err = stmtOut.Close()
			if err != nil {
				panic(err.Error())
			}
			user_free = append(user_free, 0)
			return user_free
		}
		user_free = append(user_free, one_user)
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
func DeleteChat(user_id int, my_db *sql.DB) {
	stmtIns, err := my_db.Prepare("DELETE FROM chat_buffer WHERE first_user = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtIns.Exec(user_id)
	if err != nil {
		panic(err.Error())
	}

	err = stmtIns.Close()
	if err != nil {
		panic(err.Error())
	}
}