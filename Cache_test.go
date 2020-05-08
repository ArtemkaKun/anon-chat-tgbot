package main

import (
	"reflect"
	"sort"
	"testing"
)

func Add5NewUsers() {
	for i := 0; i < 5; i++ {
		AddNewUser(int64(i))
	}
}

func Add5NewChats() {
	AddChat(1, 2)
	AddChat(3, 4)
	AddChat(5, 6)
	AddChat(7, 8)
	AddChat(9, 10)
}

func Add5Rooms() (tokens []string) {
	for i := 0; i < 5; i++ {
		tokens = append(tokens, AddRoom(int64(i)))
	}
	return
}

func TestAdd5NewUsers(t *testing.T) {
	Add5NewUsers()

	if len(GetUsersCache()) != 5 {
		t.Errorf("Expected %v, got %v", 5, len(GetUsersCache()))
	}

	for _, userStates := range GetUsersCache() {
		if userStates.IsUserSearching() || userStates.IsUserChatting() {
			t.Errorf("Expected %v, got %v", false, true)
		}
	}
}

func TestAdd5NewUsersIs7UsersExist(t *testing.T) {
	Add5NewUsers()

	var answers []bool
	for i := 0; i < 7; i++ {
		if IsUserExist(int64(i)) {
			answers = append(answers, IsUserExist(int64(i)))
		}
	}

	if len(answers) != 5 {
		t.Errorf("Expected %v, got %v", 5, len(answers))
	}
}

func TestAdd5NewUsers3UsersChangeSearchingStatus(t *testing.T) {
	Add5NewUsers()

	for i := 0; i < 3; i++ {
		ChangeUserSearchingStatus(int64(i), true)
	}

	var answers []bool
	for _, statuses := range GetUsersCache() {
		if statuses.IsUserSearching() {
			answers = append(answers, statuses.IsUserSearching())
		}
	}

	if len(answers) != 3 {
		t.Errorf("Expected %v, got %v", 3, len(answers))
	}
}

func TestAdd5NewUsers3UsersChangeChattingStatus(t *testing.T) {
	Add5NewUsers()

	for i := 0; i < 3; i++ {
		ChangeUserChattingStatus(int64(i), true)
	}

	var answers []bool
	for _, statuses := range GetUsersCache() {
		if statuses.IsUserChatting() {
			answers = append(answers, statuses.IsUserChatting())
		}
	}

	if len(answers) != 3 {
		t.Errorf("Expected %v, got %v", 3, len(answers))
	}
}

func TestAdd5NewUsers3UsersChangeSearchingStatusCheckSearchingStatus(t *testing.T) {
	Add5NewUsers()

	for i := 0; i < 3; i++ {
		ChangeUserSearchingStatus(int64(i), true)
	}

	var answers []bool
	for i := 0; i < 5; i++ {
		if CheckUserSearchingStatus(int64(i)) {
			answers = append(answers, CheckUserSearchingStatus(int64(i)))
		}
	}

	if len(answers) != 3 {
		t.Errorf("Expected %v, got %v", 3, len(answers))
	}
}

func TestAdd5NewUsers3UsersChangeChattingStatusCheckChattingStatus(t *testing.T) {
	Add5NewUsers()

	for i := 0; i < 3; i++ {
		ChangeUserChattingStatus(int64(i), true)
	}

	var answers []bool
	for i := 0; i < 5; i++ {
		if CheckUserChattingStatus(int64(i)) {
			answers = append(answers, CheckUserChattingStatus(int64(i)))
		}
	}

	if len(answers) != 3 {
		t.Errorf("Expected %v, got %v", 3, len(answers))
	}
}

func TestAdd5NewUsers3UsersChangeSearchingStatusGetSearchingUsersList(t *testing.T) {
	Add5NewUsers()

	for i := 0; i < 3; i++ {
		ChangeUserSearchingStatus(int64(i), true)
	}

	answers := []int{0, 1, 2}
	users := SearchingUsersList()

	var iusers []int
	for _, i := range users {
		iusers = append(iusers, int(i))
	}

	sort.Ints(iusers)
	if !reflect.DeepEqual(answers, iusers) {
		t.Errorf("Expected %v, got %v", answers, iusers)
	}
}

func TestAdd5Chats(t *testing.T) {
	Add5NewChats()

	if len(GetChatsCache()) != 10 {
		t.Errorf("Expected %v, got %v", 10, len(GetChatsCache()))
	}

}

func TestAdd5ChatsDelete3Chats(t *testing.T) {
	Add5NewChats()

	DeleteChat(1, 2)
	DeleteChat(5, 6)
	DeleteChat(9, 10)

	if len(GetChatsCache()) != 4 {
		t.Errorf("Expected %v, got %v", 10, len(GetChatsCache()))
	}

}

func TestAdd5Chats3TimesGetSecondUser(t *testing.T) {
	Add5NewChats()

	var users []int64

	users = append(users, GetSecondUser(1))
	users = append(users, GetSecondUser(10))
	users = append(users, GetSecondUser(5))

	answer := []int64{2, 9, 6}

	if !reflect.DeepEqual(answer, users) {
		t.Errorf("Expected %v, got %v", answer, users)
	}
}

func TestAdd5NewRooms(t *testing.T) {
	Add5Rooms()

	if len(GetRoomsCache()) != 5 {
		t.Errorf("Expected %v, got %v", 5, len(GetRoomsCache()))
	}
}

func TestAdd5NewRoomsDelete3Rooms(t *testing.T) {
	tokens := Add5Rooms()

	DeleteRoom(tokens[0])
	DeleteRoom(tokens[1])
	DeleteRoom(tokens[2])

	if len(GetRoomsCache()) != 2 {
		t.Errorf("Expected %v, got %v", 2, len(GetRoomsCache()))
	}
}

func TestAdd5NewRoomsGet3RoomsAuthor(t *testing.T) {
	tokens := Add5Rooms()

	var users []int64
	users = append(users, GetRoomAuthor(tokens[0]))
	users = append(users, GetRoomAuthor(tokens[3]))
	users = append(users, GetRoomAuthor(tokens[1]))

	answer := []int64{0, 3, 1}

	if !reflect.DeepEqual(answer, users) {
		t.Errorf("Expected %v, got %v", answer, users)
	}
}

func TestAdd5NewRoomsGet3FirstTokens(t *testing.T) {
	tokens := Add5Rooms()

	usersAnsw := []int64{0, 1, 2}
	var userTokens []string

	for _, id := range usersAnsw {
		userTokens = append(userTokens, GetRoomToken(id))
	}

	answer := []string{tokens[0], tokens[1], tokens[2]}

	if !reflect.DeepEqual(answer, userTokens) {
		t.Errorf("Expected %v, got %v", answer, userTokens)
	}
}

func TestAdd5NewRoomsCheck3UsersIsRoomAuthor(t *testing.T) {
	Add5Rooms()

	usersAnsw := []int64{0, 1, 6}
	var userAuthor []bool

	for _, id := range usersAnsw {
		userAuthor = append(userAuthor, CheckIsRoomAuthor(id))
	}

	answer := []bool{true, true, false}

	if !reflect.DeepEqual(answer, userAuthor) {
		t.Errorf("Expected %v, got %v", answer, userAuthor)
	}
}
