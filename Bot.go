package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
)

var Bot *tgbotapi.BotAPI

var chatControlKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat"),
	),
)

func init() {
	bot, err := tgbotapi.NewBotAPI("NEW BOT TOKEN HERE")
	if err != nil {
		BotInitError(err)
	}

	log.Printf("Autorised on account %s", bot.Self.UserName)

	Bot = bot
}

func BotUpdateLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		BotInitError(err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			switch update.Message.Text {
			case "New chat":
				NewChatButton(update)
				continue

			case "Leave chat":
				LeaveChatButton(update)
				continue
			}

			if FindChat(update.Message.From.ID, database, my_bot) != 0 {
				chat_id := FindChat(update.Message.From.ID, database, my_bot)
				msg := tgbotapi.NewMessage(int64(chat_id), "")

				msg.Text = update.Message.Text
				if msg.Text != "" {
					_, err := my_bot.Send(msg)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Photo != nil {
					photo := tgbotapi.NewPhotoShare(int64(chat_id), "")
					photo.FileID = (*update.Message.Photo)[2].FileID
					_, err := my_bot.Send(photo)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Voice != nil {
					voice := tgbotapi.NewVoiceShare(int64(chat_id), "")
					voice.FileID = update.Message.Voice.FileID
					_, err := my_bot.Send(voice)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Animation != nil {
					voice := tgbotapi.NewAnimationShare(int64(chat_id), "")
					voice.FileID = update.Message.Animation.FileID
					_, err := my_bot.Send(voice)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Audio != nil {
					audio := tgbotapi.NewAudioShare(int64(chat_id), "")
					audio.FileID = update.Message.Audio.FileID
					_, err := my_bot.Send(audio)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Sticker != nil {
					sticker := tgbotapi.NewStickerShare(int64(chat_id), "")
					sticker.FileID = update.Message.Sticker.FileID
					_, err := my_bot.Send(sticker)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Document != nil {
					doc := tgbotapi.NewDocumentShare(int64(chat_id), "")
					doc.FileID = update.Message.Document.FileID
					_, err := my_bot.Send(doc)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.Video != nil {
					video := tgbotapi.NewVideoShare(int64(chat_id), "")
					video.FileID = update.Message.Video.FileID
					_, err := my_bot.Send(video)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else if update.Message.VideoNote != nil {
					video_note := tgbotapi.NewVideoNoteShare(int64(chat_id), 0, "")
					video_note.FileID = update.Message.VideoNote.FileID
					video_note.Length = update.Message.VideoNote.Length
					_, err := my_bot.Send(video_note)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else {
					msg.Text = "Bot cannot send this yet! Please, contact with creator"
					_, err := my_bot.Send(msg)
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
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
			if !CheckUserReg(update.Message.From.ID, database, my_bot) {
				UserFirstStart(update.Message.From.ID, database, my_bot)
				msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers.\n\n" +
					"To start the chat, send /go_chat command or press \"New chat\" button\n\n" +
					"To leave the chat, send /leave_chat command or press \"Leave chat\" button\n\n" +
					"Bot doesn't store any personal data, so chats are fully anonymous.\n\n" +
					"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
					"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n\n" +
					"Also, check my Patreon page (https://www.patreon.com/artemkakun) if you want receive some bonuses from me :)\n"
				msg.ReplyMarkup = numericKeyboard
			} else {
				msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers.\n\n" +
					"To start the chat, send /go_chat command or press \"New chat\" button\n\n" +
					"To leave the chat, send /leave_chat command or press \"Leave chat\" button\n\n" +
					"Bot doesn't store any personal data, so chats are fully anonymous.\n\n" +
					"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
					"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n\n" +
					"Also, check my Patreon page (https://www.patreon.com/artemkakun) if you want receive some bonuses from me :)\n"
				msg.ReplyMarkup = numericKeyboard
			}
		case "go_chat":
			if CheckUserReg(update.Message.From.ID, database, my_bot) {
				if IsFree(update.Message.From.ID, database, my_bot) {
					if !IsSearch(update.Message.From.ID, database, my_bot) {
						ChangeSearch(database, update.Message.From.ID, 1, my_bot)
						msg.Text = "Search started"
					} else {
						msg.Text = "You are searching already"
					}
				} else {
					msg.Text = "You are chatting already"
				}
			} else {
				msg.Text = "You need to /start first"
			}
		case "leave_chat":
			if CheckUserReg(update.Message.From.ID, database, my_bot) {
				if FindChat(update.Message.From.ID, database, my_bot) != 0 {
					chat_id := FindChat(update.Message.From.ID, database, my_bot)
					DeleteChat(update.Message.From.ID, database, my_bot)
					ChangeState(database, update.Message.From.ID, 0, my_bot)
					msg.Text = "You leaved a chat"

					DeleteChat(chat_id, database, my_bot)
					ChangeState(database, chat_id, 0, my_bot)
					_, err := my_bot.Send(tgbotapi.NewMessage(int64(chat_id), "The stranger leave the chat"))
					if err != nil {
						ErrorCatch(err.Error(), my_bot)
					}
				} else {
					msg.Text = "You are not chatting now!"
				}
			} else {
				msg.Text = "You need to /start first"
			}
		}

		_, err := my_bot.Send(msg)
		if err != nil {
			ErrorCatch(err.Error(), my_bot)
		}
	}
}

func LeaveChatButton(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You leaved a chat")

	defer BotSendMessage(msg)

	if !CheckUserReg(update.Message.From.ID) {
		msg.Text = "You need /start first"
		return
	}

	second_user := FindSecondUserFromChat(update.Message.From.ID)

	if second_user == 0 {
		msg.Text = "You are not chatting now!"
		return
	}

	DeleteChat(update.Message.From.ID)
	ChangeUserChattingState(update.Message.From.ID, false)

	DeleteChat(second_user)
	ChangeUserChattingState(second_user, false)

	BotSendMessage(tgbotapi.NewMessage(int64(second_user), "The stranger leave the chat"))
}

func NewChatButton(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Search started")

	defer BotSendMessage(msg)

	if !CheckUserReg(update.Message.From.ID) {
		msg.Text = "You need /start first"
		return
	}

	if IsUserChatting(update.Message.From.ID) {
		msg.Text = "You are chatting already"
		return
	}

	if IsUserSearching(update.Message.From.ID) {
		msg.Text = "You are searching already"
		return
	}

	ChangeUserSearchingState(update.Message.From.ID, true)
}

func BotSendMessage(msg tgbotapi.MessageConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendMessageError(err error) {
	fmt.Println(fmt.Errorf("Send message failed: %w\n", err))
}

func BotInitError(err error) {
	fmt.Println(fmt.Errorf("Bot initialization failed: %w\n", err))
	panic(err)
}
