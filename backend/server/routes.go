package server

import (
	"net/http"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
)

func InitRoutes(r *repo.Repository) {
	Repo = r

	// Роуты для пользователей
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/updateUser", UpdateUserHandler)
	http.HandleFunc("/deleteUser", DeleteUserHandler)
	http.HandleFunc("/findUser", FindUserHandler)
	http.HandleFunc("/findUserbyID", FindUserbyIDHandler)

	// Роуты для чатов
	http.HandleFunc("/createChat", CreateChatHandler)
	http.HandleFunc("/deleteChat", DeleteChatHandler)
	http.HandleFunc("/addMember", AddMemberHandler)
	http.HandleFunc("/removeMember", RemoveMemberHandler)
	http.HandleFunc("/getChats", GetUserChatsHandler)
	http.HandleFunc("/getUsersFromChat", GetChatUsersHandler)
	http.HandleFunc("/leaveChat", LeaveChatHandler)
	http.HandleFunc("/resetUnread", UpdateLastReadHandler)
	http.HandleFunc("/unreadCounts", GetUnreadCountsHandler)

	// Роуты для сообщений
	http.HandleFunc("/deleteMessage", DeleteMessageHandler)
	http.HandleFunc("/getMessage", GetMessagesHandler)
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/ws", WSHandler(hub))
}
