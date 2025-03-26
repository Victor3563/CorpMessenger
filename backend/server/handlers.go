// Organization conection(hand) in different situation like reg, update, delete and auth
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
	"github.com/Victor3563/CorpMessenger/root"
)

var Repo *repo.Repository

// regist new user by User.Create and write it to server
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server regist comand get")
	req, err := root.ParserandValidByNameandMethod(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := Repo.CreateUser(req.Name, req.Password, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Auth user by User.Auth in server
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Auth comand get")
	req, err := root.ParserandValidByNameandMethod(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := Repo.AuthenticateUser(req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Update user by User.Update and write it to server
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Update comand get")
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req repo.User
	req, err := root.ParserandValidByName(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := Repo.UpdateUser(req.ID, req.Name, req.Password, req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// Delete user by User.Delete and commit change on server
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Delete comand get")
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	req, err := root.ParserandValidByID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := Repo.DeleteUser(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}

// FindUserHandler возвращает список чатов для пользователя
func FindUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Find users comand get")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("username")
	if name == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	chats, err := Repo.FindUser(name)
	if err != nil {
		http.Error(w, "Error fetching chats: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, "Error formating answer", http.StatusInternalServerError)
	}

}

// FindUserbyIDHandler возвращает информацию о пользователе
func FindUserbyIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Find users comand get")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	user, err := Repo.FindUserbyID(userID)
	if err != nil {
		http.Error(w, "Error fetching chats: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error formating answer", http.StatusInternalServerError)
	}

}
