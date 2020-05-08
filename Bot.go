package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
	"strings"
)

var Bot *tgbotapi.BotAPI

var chatControlKeyboardUS = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Create a secret room"),
		tgbotapi.NewKeyboardButton("Join a secret room"),
	),
)

var chatControlKeyboardRU = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Начать чат\n"),
		tgbotapi.NewKeyboardButton("Покинуть чат\n"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создать секретную комнату\n"),
		tgbotapi.NewKeyboardButton("Войти в секретную комнату\n"),
	),
)

func init() {
	var err error

	Bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Autorised on account %s", Bot.Self.UserName)
}

func BotUpdateLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			switch update.Message.Text {
			case "New chat", "Начать чат":
				NewChatButton(update)
				continue

			case "Leave chat", "Покинуть чат":
				LeaveChatButton(update)
				continue

			case "Create a secret room", "Создать секретную комнату":
				CreateSecretRoom(update)
				continue

			case "Join a secret room", "Войти в секретную комнату":
				msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Now provide the secret token. Please, type word 'token', space and provide token")

				if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
					msg.Text = "Теперь введите секретный токен. Пожалуйста, напишите слово 'token', пробел, и введите токен"
				}

				BotSendMessage(msg)
				continue
			}

			if strings.Contains(update.Message.Text, "token") {
				form := strings.Fields(update.Message.Text)
				if len(form) < 2 {
					msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Token was provided wrongly")

					if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
						msg.Text = "Токен был введен не правильно"
					}

					BotSendMessage(msg)
					continue
				}

				token := strings.Fields(update.Message.Text)[1]

				if GetRoomAuthor(token) == 0 {
					msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Cannot find a secret room with that token")

					if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
						msg.Text = "Не могу найти комнату с таким токеном"
					}

					BotSendMessage(msg)
					continue
				}

				JoinSecretRoom(update, token)
			}

			if !CheckUserChattingStatus(int64(update.Message.From.ID)) {
				continue
			}

			SendMessageToAnotherUser(update)
			continue
		}

		switch update.Message.Command() {
		case "start":
			StartCommand(update)
		}
	}
}

func NewChatButton(update tgbotapi.Update) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	if CheckUserSearchingStatus(int64(update.Message.From.ID)) {
		AlreadySearching(int64(update.Message.From.ID), update)
		return
	}

	if CheckUserChattingStatus(int64(update.Message.From.ID)) {
		AlreadyChatting(int64(update.Message.From.ID), update)
		return
	}

	if CheckIsRoomAuthor(int64(update.Message.From.ID)) {
		RoomAuthor(int64(update.Message.From.ID), update)
		return
	}

	ChangeUserSearchingStatus(int64(update.Message.From.ID), true)

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Search started")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Поиск начался")
	}

	BotSendMessage(msg)
}

func AlreadySearching(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are searching already")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже ищите чат"
	}

	BotSendMessage(msg)
}

func AlreadyChatting(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are chatting already")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже в чате"
	}

	BotSendMessage(msg)
}

func RoomAuthor(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are the secret room author. Wait for second people or leave the room")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы автор секретной комнаты. Подождите второго человека или выйдите из комнаты"
	}

	BotSendMessage(msg)
}

func LeaveChatButton(update tgbotapi.Update) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	if !CheckUserChattingStatus(int64(update.Message.From.ID)) {
		NotChatting(int64(update.Message.From.ID), update)
		return
	}

	if CheckIsRoomAuthor(int64(update.Message.From.ID)) {
		DeleteRoom(GetRoomToken(int64(update.Message.From.ID)))

		msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You delete a secret room")

		if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
			msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Вы удалили секретную комнату")
		}

		BotSendMessage(msg)
		return
	}

	secondUser := GetSecondUser(int64(update.Message.From.ID))

	DeleteChat(int64(update.Message.From.ID), secondUser)

	ChangeUserChattingStatus(int64(update.Message.From.ID), false)
	ChangeUserChattingStatus(secondUser, false)

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You leaved a chat")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Вы покинули чат")
	}

	BotSendMessage(msg)
	BotSendMessage(tgbotapi.NewMessage(secondUser, "The stranger leave the chat"))
}

func NotChatting(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are not chatting now!")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы не находитесь в чате"
	}

	BotSendMessage(msg)
}

func CreateSecretRoom(update tgbotapi.Update) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	if CheckUserSearchingStatus(int64(update.Message.From.ID)) {
		AlreadySearching(int64(update.Message.From.ID), update)
		return
	}

	if CheckUserChattingStatus(int64(update.Message.From.ID)) {
		AlreadyChatting(int64(update.Message.From.ID), update)
		return
	}

	if CheckIsRoomAuthor(int64(update.Message.From.ID)) {
		RoomAuthor(int64(update.Message.From.ID), update)
		return
	}

	token := AddRoom(int64(update.Message.From.ID))

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), fmt.Sprintf("You created a secret room. Your secret token %v. Share it to another person.", token))

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = fmt.Sprintf("Вы создали секретную комнату. Ваш секретный токен %v. Передайте его другому человеку, что бы он мог подключиться к комнате.", token)
	}

	BotSendMessage(msg)
}

func JoinSecretRoom(update tgbotapi.Update, token string) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	if CheckUserSearchingStatus(int64(update.Message.From.ID)) {
		AlreadySearching(int64(update.Message.From.ID), update)
		return
	}

	if CheckUserChattingStatus(int64(update.Message.From.ID)) {
		AlreadyChatting(int64(update.Message.From.ID), update)
		return
	}

	if CheckIsRoomAuthor(int64(update.Message.From.ID)) {
		RoomAuthor(int64(update.Message.From.ID), update)
		return
	}

	roomAuthor := GetRoomAuthor(token)
	DeleteRoom(token)

	ChangeUserChattingStatus(int64(update.Message.From.ID), true)
	ChangeUserChattingStatus(roomAuthor, true)

	AddChat(int64(update.Message.From.ID), roomAuthor)

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You joined a secret room")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы присоединились к секретной комнате"
	}

	BotSendMessage(msg)
	BotSendMessage(tgbotapi.NewMessage(roomAuthor, "Another person joined a secret room"))
}

func StartCommand(update tgbotapi.Update) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}
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

	BotSendMessage(msg)
}

func SendMessageToAnotherUser(update tgbotapi.Update) {
	secondUser := GetSecondUser(int64(update.Message.From.ID))

	msg := tgbotapi.NewMessage(secondUser, update.Message.Text)

	if msg.Text != "" {
		BotSendMessage(msg)
		return
	}

	if update.Message.Photo != nil {
		photo := tgbotapi.NewPhotoShare(secondUser, "")
		photo.FileID = (*update.Message.Photo)[1].FileID

		BotSendPhoto(photo)
		return
	}

	if update.Message.Voice != nil {
		voice := tgbotapi.NewVoiceShare(secondUser, "")
		voice.FileID = update.Message.Voice.FileID

		BotSendVoice(voice)
		return
	}

	if update.Message.Animation != nil {
		gif := tgbotapi.NewAnimationShare(secondUser, "")
		gif.FileID = update.Message.Animation.FileID

		BotSendGif(gif)
		return
	}

	if update.Message.Audio != nil {
		audio := tgbotapi.NewAudioShare(secondUser, "")
		audio.FileID = update.Message.Audio.FileID

		BotSendAudio(audio)
		return
	}

	if update.Message.Sticker != nil {
		sticker := tgbotapi.NewStickerShare(secondUser, "")
		sticker.FileID = update.Message.Sticker.FileID

		BotSendSticker(sticker)
		return
	}

	if update.Message.Document != nil {
		doc := tgbotapi.NewDocumentShare(secondUser, "")
		doc.FileID = update.Message.Document.FileID

		BotSendDoc(doc)
		return
	}

	if update.Message.Video != nil {
		video := tgbotapi.NewVideoShare(secondUser, "")
		video.FileID = update.Message.Video.FileID

		BotSendVideo(video)
		return
	}

	if update.Message.VideoNote != nil {
		videoNote := tgbotapi.NewVideoNoteShare(secondUser, 0, "")
		videoNote.FileID = update.Message.VideoNote.FileID
		videoNote.Length = update.Message.VideoNote.Length

		BotSendVideoNote(videoNote)
		return
	}

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Бот не может это отправить. Пожалуйста, свяжитесь с администрацией"
	} else {
		msg.Text = "Bot cannot send this yet! Please, contact with creator"
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
