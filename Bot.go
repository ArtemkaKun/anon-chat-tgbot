package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
)

var Bot *tgbotapi.BotAPI

var chatControlKeyboardUS = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat"),
	),
)

var chatControlKeyboardRU = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Начать чат"),
		tgbotapi.NewKeyboardButton("Покинуть чат"),
	),
)

func init() {
	bot, err := tgbotapi.NewBotAPI("1022122500:AAFy8sDJFUlgw0e7JURelghBPv_is5kG7ck")
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

			if FindSecondUserFromChat(update.Message.From.ID) == 0 {
				continue
			}

			SendMessageToAnotherUser(update)
			continue
		}

		switch update.Message.Command() {
		case "start":
			StartCommand(update)

		case "go_chat":
			NewChatButton(update)

		case "leave_chat":
			LeaveChatButton(update)
		}
	}
}

func StartCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Привет, это Freenon чат - анонимный чат, где ты можешь высказывать свои мысли без последствий\n\n" +
			"Чтобы начать чат с незнакомцем, введи команду /go_chat или нажми кнопку \"Начать чат\"\n\n" +
			"Чтобы покинуть чат, введи команду /leave_chat или нажми кнопку \"Покинуть чат\"\n\n" +
			"Бот не сохраняет данные о пользователях, так что твои личные данные в безопасности.\n\n" +
			"Если ты хочешь посмотреть, как работает бот - вот мое видео (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
			"Если у тебя есть вопросы или предложения, пожалуйста, свяжись со мной, @YUART\n\n" +
			"Еще можешь зайти на мою страницу Parteon (https://www.patreon.com/artemkakun) - мне будет приятно :)\n"
		msg.ReplyMarkup = chatControlKeyboardRU
	} else {
		msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers.\n\n" +
			"To start the chat, send /go_chat command or press \"New chat\" button\n\n" +
			"To leave the chat, send /leave_chat command or press \"Leave chat\" button\n\n" +
			"Bot doesn't store any personal data, so chats are fully anonymous.\n\n" +
			"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
			"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n\n" +
			"Also, check my Patreon page (https://www.patreon.com/artemkakun) if you want receive some bonuses from me :)\n"
		msg.ReplyMarkup = chatControlKeyboardUS
	}

	defer BotSendMessage(msg)

	if !CheckUserReg(update.Message.From.ID) {
		UserFirstStart(update.Message.From.ID)
	}
}

func SendMessageToAnotherUser(update tgbotapi.Update) {
	second_user := FindSecondUserFromChat(update.Message.From.ID)

	msg := tgbotapi.NewMessage(int64(second_user), update.Message.Text)

	if msg.Text != "" {
		BotSendMessage(msg)
		return
	}

	if update.Message.Photo != nil {
		photo := tgbotapi.NewPhotoShare(int64(second_user), "")
		photo.FileID = (*update.Message.Photo)[1].FileID

		BotSendPhoto(photo)
		return
	}

	if update.Message.Voice != nil {
		voice := tgbotapi.NewVoiceShare(int64(second_user), "")
		voice.FileID = update.Message.Voice.FileID

		BotSendVoice(voice)
		return
	}

	if update.Message.Animation != nil {
		gif := tgbotapi.NewAnimationShare(int64(second_user), "")
		gif.FileID = update.Message.Animation.FileID

		BotSendGif(gif)
		return
	}

	if update.Message.Audio != nil {
		audio := tgbotapi.NewAudioShare(int64(second_user), "")
		audio.FileID = update.Message.Audio.FileID

		BotSendAudio(audio)
		return
	}

	if update.Message.Sticker != nil {
		sticker := tgbotapi.NewStickerShare(int64(second_user), "")
		sticker.FileID = update.Message.Sticker.FileID

		BotSendSticker(sticker)
		return
	}

	if update.Message.Document != nil {
		doc := tgbotapi.NewDocumentShare(int64(second_user), "")
		doc.FileID = update.Message.Document.FileID

		BotSendDoc(doc)
		return
	}

	if update.Message.Video != nil {
		video := tgbotapi.NewVideoShare(int64(second_user), "")
		video.FileID = update.Message.Video.FileID

		BotSendVideo(video)
		return
	}

	if update.Message.VideoNote != nil {
		video_note := tgbotapi.NewVideoNoteShare(int64(second_user), 0, "")
		video_note.FileID = update.Message.VideoNote.FileID
		video_note.Length = update.Message.VideoNote.Length

		BotSendVideoNote(video_note)
		return
	}

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Бот не может это отправить. Пожалуйста, свяжитесь с администрацией"
	} else {
		msg.Text = "Bot cannot send this yet! Please, contact with creator"
	}
	BotSendMessage(msg)
}

func LeaveChatButton(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You leaved a chat")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Вы покинули чат")
	}

	if !CheckUserReg(update.Message.From.ID) {
		BadRegistration(msg, update)
		return
	}

	second_user := FindSecondUserFromChat(update.Message.From.ID)

	if second_user == 0 {
		NotChatting(msg, update)
		return
	}

	DeleteChat(update.Message.From.ID)
	ChangeUserChattingState(update.Message.From.ID, false)

	DeleteChat(second_user)
	ChangeUserChattingState(second_user, false)

	BotSendMessage(msg)
	BotSendMessage(tgbotapi.NewMessage(int64(second_user), "The stranger leave the chat"))
}

func NewChatButton(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Search started")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Поиск начался")
	}

	if !CheckUserReg(update.Message.From.ID) {
		BadRegistration(msg, update)
		return
	}

	if IsUserChatting(update.Message.From.ID) {
		AlreadyChatting(msg, update)
		return
	}

	if IsUserSearching(update.Message.From.ID) {
		AlreadySearching(msg, update)
		return
	}

	ChangeUserSearchingState(update.Message.From.ID, true)
	BotSendMessage(msg)
}

func AlreadySearching(msg tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "You are searching already"
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже ищите чат"
	}

	BotSendMessage(msg)
}

func AlreadyChatting(msg tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "You are chatting already"
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже в чате"
	}

	BotSendMessage(msg)
}

func BadRegistration(msg tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "You need /start first"
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Используйте /start сначала"
	}

	BotSendMessage(msg)
}

func NotChatting(msg tgbotapi.MessageConfig, update tgbotapi.Update) {
	msg.Text = "You are not chatting now!"
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы не находитесь в чате"
	}

	BotSendMessage(msg)
}

func BotSendMessage(msg tgbotapi.MessageConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendPhoto(msg tgbotapi.PhotoConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendVoice(msg tgbotapi.VoiceConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendGif(msg tgbotapi.AnimationConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendAudio(msg tgbotapi.AudioConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendSticker(msg tgbotapi.StickerConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendDoc(msg tgbotapi.DocumentConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendVideo(msg tgbotapi.VideoConfig) {
	_, err := Bot.Send(msg)
	if err != nil {
		BotSendMessageError(err)
	}
}

func BotSendVideoNote(msg tgbotapi.VideoNoteConfig) {
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
