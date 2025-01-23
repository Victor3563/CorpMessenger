package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Victor3563/CorpMessenger/server"
)

func main() {
	http.HandleFunc("/register", server.RegisterHandler)
	http.HandleFunc("/auth", server.AuthHandler)
	http.HandleFunc("/update", server.UpdateUserHandler)
	http.HandleFunc("/delete", server.DeleteUserHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
