package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
	"github.com/Victor3563/CorpMessenger/server"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://user:password@localhost:5432/messenger?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация репозитория и передача его в глобальную переменную обработчиков
	repoInstance := repo.NewRepository(db)
	server.Repo = repoInstance

	// Роуты для пользователей
	http.HandleFunc("/register", server.RegisterHandler)
	http.HandleFunc("/auth", server.AuthHandler)
	http.HandleFunc("/updateUser", server.UpdateUserHandler)
	http.HandleFunc("/deleteUser", server.DeleteUserHandler)

	// Роуты для чатов
	http.HandleFunc("/createChat", server.CreateChatHandler)
	http.HandleFunc("/deleteChat", server.DeleteChatHandler)
	http.HandleFunc("/addMember", server.AddMemberHandler)
	http.HandleFunc("/removeMember", server.RemoveMemberHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
