//REalesation User config + funk like Create, Auth, Update and delete

package repo

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

var (
	UserData = make(map[string]User)
	mu       sync.Mutex
)

// find user in User Data by name
func findUserByUsername(username string) (User, bool) {
	for _, user := range UserData {
		if user.Name == username {
			return user, true
		}
	}
	return User{}, false
}

// Create new user
func CreateUser(username, password, email string) (User, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if username already exists
	if _, exists := findUserByUsername(username); exists {
		return User{}, errors.New("username already exists")
	}

	id := uuid.New().String()
	user := User{
		ID:       id,
		Name:     username,
		Password: password,
		Email:    email,
	}
	UserData[id] = user
	return user, nil
}

// AuthenticateUser checks if a username and password are correct
func AuthenticateUser(username, password string) (User, error) {
	mu.Lock()
	defer mu.Unlock()

	user, exists := findUserByUsername(username)
	if !exists || user.Password != password {
		return User{}, errors.New("invalid username or password")
	}
	return user, nil
}

// UpdateUser updates user attributes in the store
func UpdateUser(id string, newUsername, newPassword, newEmail string) error {
	mu.Lock()
	defer mu.Unlock()

	user, exists := UserData[id]
	if !exists {
		return errors.New("user not found")
	}

	if newUsername != "" {
		user.Name = newUsername
	}
	if newPassword != "" {
		user.Password = newPassword
	}
	if newEmail != "" {
		user.Email = newEmail
	}

	UserData[id] = user
	return nil
}

// DeleteUser removes a user from the Data
func DeleteUser(id string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := UserData[id]; !exists {
		return errors.New("user not found")
	}

	delete(UserData, id)
	return nil
}
