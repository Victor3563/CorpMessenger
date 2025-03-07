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

	// Роуты для чатов
	http.HandleFunc("/createChat", CreateChatHandler)
	http.HandleFunc("/deleteChat", DeleteChatHandler)
	http.HandleFunc("/addMember", AddMemberHandler)
	http.HandleFunc("/removeMember", RemoveMemberHandler)
}
