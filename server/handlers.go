// Organization conection(hand) in different situation like reg, update, delete and auth
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
	"github.com/Victor3563/CorpMessenger/root"
)

var Repo *repo.Repository

// regist new user by User.Create and write it to server
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server regist comand get")
	req, err := root.ParserandValidByName(r)
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
	req, err := root.ParserandValidByName(r)
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

	req, err := root.ParserandValidByName(r)
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
