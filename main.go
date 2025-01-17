package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/auth", AuthHandler)
	http.HandleFunc("/update", UpdateUserHandler)
	http.HandleFunc("/delete", DeleteUserHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
