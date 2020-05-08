package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

type UserStatuses struct {
	SearchingStatus bool
	ChattingStatus  bool
}

func (user *UserStatuses) SetSearchingStatus(state bool) {
	user.SearchingStatus = state
}

func (user *UserStatuses) SetChattingStatus(state bool) {
	user.ChattingStatus = state
}

func (user *UserStatuses) IsUserSearching() bool {
	return user.SearchingStatus
}

func (user *UserStatuses) IsUserChatting() bool {
	return user.ChattingStatus
}

var Users map[int64]*UserStatuses
var Chats map[int64]int64
var Rooms map[string]int64

func init() {
	Users = make(map[int64]*UserStatuses)
	for key, val := range GetUsersFromDB() {
		Users[key] = val
	}

	Chats = make(map[int64]int64)
	for key, val := range GetChatsFromDB() {
		Chats[key] = val
	}

	Rooms = make(map[string]int64)
	for key, val := range GetRoomsFromDB() {
		Rooms[key] = val
	}
}

func GetUsersCache() map[int64]*UserStatuses {
	return Users
}

func GetChatsCache() map[int64]int64 {
	return Chats
}

func GetRoomsCache() map[string]int64 {
	return Rooms
}

func IsUserExist(user int64) bool {
	_, exist := Users[user]
	return exist
}

func AddNewUser(user int64) {
	Users[user] = new(UserStatuses)
	Users[user].SetSearchingStatus(false)
	Users[user].SetChattingStatus(false)
}

func ChangeUserSearchingStatus(user int64, status bool) {
	Users[user].SetSearchingStatus(status)
}

func ChangeUserChattingStatus(user int64, status bool) {
	Users[user].SetChattingStatus(status)
}

func CheckUserSearchingStatus(user int64) bool {
	return Users[user].IsUserSearching()
}

func CheckUserChattingStatus(user int64) bool {
	return Users[user].IsUserChatting()
}

func SearchingUsersList() (users []int64) {
	for userId, userStatus := range Users {
		if userStatus.IsUserSearching() {
			users = append(users, userId)
		}
	}

	return
}

func AddChat(firstUser, secondUser int64) {
	Chats[firstUser] = secondUser
	Chats[secondUser] = firstUser
}

func DeleteChat(firstUser, secondUser int64) {
	delete(Chats, firstUser)
	delete(Chats, secondUser)
}

func GetSecondUser(firstUser int64) int64 {
	return Chats[firstUser]
}

func AddRoom(authorUser int64) (token string) {
	token = CreateToken(authorUser)
	Rooms[token] = authorUser
	return
}

func DeleteRoom(token string) {
	delete(Rooms, token)
}

func GetRoomAuthor(token string) int64 {
	return Rooms[token]
}

func GetRoomToken(user int64) string {
	for token, roomAuthor := range Rooms {
		if roomAuthor == user {
			return token
		}
	}

	return ""
}

func CheckIsRoomAuthor(user int64) bool {
	for _, roomAuthor := range Rooms {
		if roomAuthor == user {
			return true
		}
	}

	return false
}

func CreateToken(id int64) string {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%v%v", id, time.Now())))
	return hex.EncodeToString(hasher.Sum(nil))
}
