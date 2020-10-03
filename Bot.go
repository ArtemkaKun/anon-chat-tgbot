package main

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"log"
	"os"
	"strings"
)

const (
	MIN_WORDS_IN_TOKEN_FORM = 2
	TOKEN_POSITION_IN_FORM  = 1
)

var Bot *tgbotapi.BotAPI

var chatControlKeyboardUS = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("New chat"),
		tgbotapi.NewKeyboardButton("Leave chat/room"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Create a secret room"),
		tgbotapi.NewKeyboardButton("Join a secret room"),
	),
)

var chatControlKeyboardRU = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Начать чат"),
		tgbotapi.NewKeyboardButton("Покинуть чат/комнату"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создать секретную комнату"),
		tgbotapi.NewKeyboardButton("Войти в секретную комнату"),
	),
)

var chatControlKeyboardMD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Incepe conversatia"),
		tgbotapi.NewKeyboardButton("Părăsește chatul/camera"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Creați o cameră secretă"),
		tgbotapi.NewKeyboardButton("Intră în camera secretă"),
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

	_, err := Bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("PUBLIC_URL") + ":8443" + "/" + os.Getenv("BOT_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}

	updates := Bot.ListenForWebhook("/" + os.Getenv("BOT_TOKEN"))

	//updates, err = Bot.GetUpdatesChan(u)
	//if err != nil {
	//	log.Panic(err)
	//}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			if !IsUserExist(int64(update.Message.From.ID)) {
				AddNewUser(int64(update.Message.From.ID))
			}

			switch update.Message.Text {
			case "New chat", "Начать чат", "Incepe conversatia":
				NewChatButton(update)
				continue
			case "Leave chat/room", "Покинуть чат/комнату", "Părăsește chatul/camera":
				LeaveChatButton(update)
				continue
			case "Create a secret room", "Создать секретную комнату", "Creați o cameră secretă":
				CreateSecretRoom(update)
				continue
			case "Join a secret room", "Войти в секретную комнату", "Intră în camera secretă":
				JoinRoomTokenMessage(update)
				continue
			}

			if strings.Contains(update.Message.Text, "token") {
				formData := strings.Fields(update.Message.Text)

				if len(formData) < MIN_WORDS_IN_TOKEN_FORM {
					InvalidTokenFormError(update)
					continue
				}

				token := formData[TOKEN_POSITION_IN_FORM]

				if GetRoomAuthor(token) == 0 {
					InvalidTokenError(update)
					continue
				}

				JoinSecretRoom(update, token)
				continue
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
	if usersChecks(update) {
		return
	}

	userID := int64(update.Message.From.ID)
	ChangeUserSearchingStatus(true, userID)

	msg := tgbotapi.NewMessage(userID, "Search started")
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(userID, "Поиск начался")
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg = tgbotapi.NewMessage(userID, "Căutarea a început")
	}

	BotSendMessage(msg)
}

func LeaveChatButton(update tgbotapi.Update) {
	userID := int64(update.Message.From.ID)

	if !IsUserExist(userID) {
		AddNewUser(userID)
	}

	if CheckIsRoomAuthor(userID) {
		DeleteRoom(GetRoomToken(userID))
		DeleteRoomMessage(update, userID)
		return
	}

	if !CheckUserChattingStatus(userID) {
		NotChattingError(userID, update)
		return
	}

	RemoveChat(update, userID)
}

func RemoveChat(update tgbotapi.Update, userID int64) {
	secondUser := GetSecondUser(userID)

	DeleteChat(userID, secondUser)

	ChangeUserChattingStatus(false, userID, secondUser)

	LeaveChatMessage(update, userID, secondUser)
}

func CreateSecretRoom(update tgbotapi.Update) {
	if usersChecks(update) {
		return
	}

	token := AddRoom(int64(update.Message.From.ID))

	CreateRoomMessage(update, token)
}

func JoinSecretRoom(update tgbotapi.Update, token string) {
	if usersChecks(update) {
		return
	}

	roomAuthor := GetRoomAuthor(token)
	DeleteRoom(token)

	ChangeUserChattingStatus(true, int64(update.Message.From.ID), roomAuthor)

	AddChat(int64(update.Message.From.ID), roomAuthor)

	JoinRoomMessage(update, roomAuthor)
}

func usersChecks(update tgbotapi.Update) bool {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	if CheckUserSearchingStatus(int64(update.Message.From.ID)) {
		AlreadySearchingError(int64(update.Message.From.ID), update)
		return true
	}

	if CheckUserChattingStatus(int64(update.Message.From.ID)) {
		AlreadyChattingError(int64(update.Message.From.ID), update)
		return true
	}

	if CheckIsRoomAuthor(int64(update.Message.From.ID)) {
		RoomAuthorError(int64(update.Message.From.ID), update)
		return true
	}
	return false
}

func StartCommand(update tgbotapi.Update) {
	if !IsUserExist(int64(update.Message.From.ID)) {
		AddNewUser(int64(update.Message.From.ID))
	}

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Привет, это Freenon чат - анонимный чат, где ты можешь высказывать свои мысли без последствий.\n\n" +
			"Чтобы начать чат с незнакомцем, нажми кнопку \"Начать чат\"\n\n" +
			"Чтобы покинуть чат, нажми кнопку \"Покинуть чат/комнату\"\n\n" +
			"Бот не сохраняет данные о пользователях, так что твои личные данные в безопасности.\n\n" +
			"Если ты хочешь посмотреть, как работает бот - вот мое видео (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
			"Если у тебя есть вопросы или предложения, пожалуйста, свяжись со мной, @YUART\n\n"
		msg.ReplyMarkup = chatControlKeyboardRU
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Bună, acesta este chatul AnonChatMoldova - un chat anonim în care îți poți exprima gândurile fără consecințe.\n\n" +
			"Pentru a începe un chat cu un străin, dă click pe butonul \"Incepe conversatia\".\n\n" +
			"Pentru a părăsi chatul, dă click pe butonul \"Părăsește chatul/camera\".\n\n" +
			"Botul nu salvează date despre utilizatori astfel încât informațiile dvs. personale să fie în siguranță.\n\n" +
			"Dacă aveți întrebări sau sugestii, vă rugăm să mă contactați @OWNERMD"
		msg.ReplyMarkup = chatControlKeyboardMD
	} else {
		msg.Text = "Hello, this is Freedom chat, where you can freely express your minds and talk with other strangers.\n\n" +
			"To start the chat, press \"New chat\" button\n\n" +
			"To leave the chat, press \"Leave chat/room\" button\n\n" +
			"Bot doesn't store any personal data, so chats are fully anonymous.\n\n" +
			"If You want to check how the bot works - check my video (https://www.youtube.com/watch?v=drtAdOByW54&t=1s)\n\n" +
			"If You have some questions or suggestions, please, feel free to contact with me, @YUART\n\n"
		msg.ReplyMarkup = chatControlKeyboardUS
	}

	BotSendMessage(msg)
}

func SendMessageToAnotherUser(update tgbotapi.Update) {
	secondUser := GetSecondUser(int64(update.Message.From.ID))

	var msg tgbotapi.Chattable

	if update.Message.Text != "" {
		msg = tgbotapi.NewMessage(secondUser, update.Message.Text)
	}

	if update.Message.Photo != nil {
		photo := tgbotapi.NewPhotoShare(secondUser, "")
		photo.FileID = (*update.Message.Photo)[1].FileID

		msg = photo
	}

	if update.Message.Voice != nil {
		voice := tgbotapi.NewVoiceShare(secondUser, "")
		voice.FileID = update.Message.Voice.FileID

		msg = voice
	}

	if update.Message.Animation != nil {
		gif := tgbotapi.NewAnimationShare(secondUser, "")
		gif.FileID = update.Message.Animation.FileID

		msg = gif
	}

	if update.Message.Audio != nil {
		audio := tgbotapi.NewAudioShare(secondUser, "")
		audio.FileID = update.Message.Audio.FileID

		msg = audio
	}

	if update.Message.Sticker != nil {
		sticker := tgbotapi.NewStickerShare(secondUser, "")
		sticker.FileID = update.Message.Sticker.FileID

		msg = sticker
	}

	if update.Message.Document != nil {
		doc := tgbotapi.NewDocumentShare(secondUser, "")
		doc.FileID = update.Message.Document.FileID

		msg = doc
	}

	if update.Message.Video != nil {
		video := tgbotapi.NewVideoShare(secondUser, "")
		video.FileID = update.Message.Video.FileID

		msg = video
	}

	if update.Message.VideoNote != nil {
		videoNote := tgbotapi.NewVideoNoteShare(secondUser, 0, "")
		videoNote.FileID = update.Message.VideoNote.FileID
		videoNote.Length = update.Message.VideoNote.Length

		msg = videoNote
	}

	if msg == nil {
		if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
			msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Бот не может это отправить. Пожалуйста, свяжитесь с администрацией")
		} else if update.Message.From.LanguageCode == "ro-MD" {
			msg = tgbotapi.NewMessage(int64(update.Message.From.ID),"Botul nu poate trimite asta. Vă rugăm să contactați administrația")
		} else {
			msg = tgbotapi.NewMessage(int64(update.Message.From.ID), "Bot cannot send this yet! Please, contact with creator")
		}
	}

	BotSendMessage(msg)
}

func BotSendMessage(message tgbotapi.Chattable) {
	var err error

	_, err = Bot.Send(message)

	if err != nil {
		BotSendMessageError(err)
	}
}

func JoinRoomTokenMessage(update tgbotapi.Update) {
	if usersChecks(update) {
		return
	}

	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Now provide the secret token. Please, type word 'token', space and provide token")
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Теперь введите секретный токен. Пожалуйста, напишите слово 'token', пробел, и введите токен"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Acum introduceți simbolul secret. Vă rugăm să scrieți cuvântul „token”, un spațiu și introduceți simbolul"
	}

		BotSendMessage(msg)
}

func LeaveChatMessage(update tgbotapi.Update, userID int64, secondUser int64) {
	msg := tgbotapi.NewMessage(userID, "You leaved a chat")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(userID, "Вы покинули чат")
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg = tgbotapi.NewMessage(userID, "Ai părăsit chatul")
	}

	BotSendMessage(msg)
	BotSendMessage(tgbotapi.NewMessage(secondUser, "The stranger leave the chat"))
}

func DeleteRoomMessage(update tgbotapi.Update, userID int64) {
	msg := tgbotapi.NewMessage(userID, "You delete a secret room")
	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg = tgbotapi.NewMessage(userID, "Вы удалили секретную комнату")
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg = tgbotapi.NewMessage(userID, "Ați șters camera secretă")
	}

	BotSendMessage(msg)
}

func CreateRoomMessage(update tgbotapi.Update, token string) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), fmt.Sprintf("You created a secret room. Your secret token %v. Share it to another person.", token))

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = fmt.Sprintf("Вы создали секретную комнату. Ваш секретный токен %v. Передайте его другому человеку, что бы он мог подключиться к комнате.", token)
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = fmt.Sprintf("Ați creat o cameră secretă. Token-ul tău secret %v. Dă-i altei persoane, astfel încât să se poată conecta la cameră.", token)
	}

	BotSendMessage(msg)
}

func JoinRoomMessage(update tgbotapi.Update, roomAuthor int64) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "You joined a secret room")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы присоединились к секретной комнате"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Te-ai alăturat camerei secrete"
	}

	BotSendMessage(msg)
	BotSendMessage(tgbotapi.NewMessage(roomAuthor, "Another person joined a secret room"))
}

func BotSendMessageError(err error) {
	fmt.Println(fmt.Errorf("Send message failed: %w\n", err))
}

func InvalidTokenError(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Cannot find a secret room with that token")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Не могу найти комнату с таким токеном"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Nu găsesc o cameră cu un astfel de Token"
	}

	BotSendMessage(msg)
}

func InvalidTokenFormError(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(int64(update.Message.From.ID), "Token was provided wrongly")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Токен был введен не правильно"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Token-ul a fost introdus incorect"
	}

	BotSendMessage(msg)
}

func AlreadySearchingError(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are searching already")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже ищите чат"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Căutați deja un chat"
	}

	BotSendMessage(msg)
}

func AlreadyChattingError(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are chatting already")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы уже в чате"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Sunteti deja in chat"
	}

	BotSendMessage(msg)
}

func RoomAuthorError(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are the secret room author. Wait for second people or leave the room")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы автор секретной комнаты. Подождите второго человека или выйдите из комнаты"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Ești autorul camerei secrete. Așteptați a doua persoană sau părăsiți camera"
	}

	BotSendMessage(msg)
}

func NotChattingError(user int64, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(user, "You are not chatting now!")

	if update.Message.From.LanguageCode == "ru" || update.Message.From.LanguageCode == "ua" {
		msg.Text = "Вы не находитесь в чате"
	} else if update.Message.From.LanguageCode == "ro-MD" {
		msg.Text = "Nu va aflati in niciun chat"
	}

	BotSendMessage(msg)
}
