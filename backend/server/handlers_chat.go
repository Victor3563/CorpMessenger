package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Victor3563/CorpMessenger/root"
)

func CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server create chat comand get")
	req, err := root.ParserConversationAdd(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conv, err := Repo.CreateConversation(req.Type, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	creatorID := req.CreatorID
	if err := Repo.AddMemberToConversation(conv.ID, creatorID, "admin"); err != nil {
		http.Error(w, "Chat created but error adding creator: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conv)
}
func LeaveChatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Leave chat command received")
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ConversationID int `json:"conversation_id"`
		UserID         int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Пользователь удаляет себя; здесь не требуется проверка админа
	if err := Repo.RemoveMemberFromConversation(req.ConversationID, req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Left chat successfully"))
}
func DeleteChatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server Delete chat comand get")
	req, err := root.ParserConversationDelete(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := Repo.DeleteConversation(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Conversation deleted successfully"))
}

func AddMemberHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server add to chat comand get")
	req, err := root.ParserAddToConversation(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := Repo.AddMemberToConversation(req.ConversationID, req.UserID, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Member added successfully"))
}

func RemoveMemberHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("remove member comand get")
	req, err := root.ParserRemoveFromConversation(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := Repo.RemoveMemberFromConversation(req.ConversationID, req.UserID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Member removed successfully"))

}

// GetUserChatsHandler возвращает список чатов для пользователя
func GetUserChatsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get chats comand get")
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
	chats, err := Repo.GetUserChats(userID)
	if err != nil {
		http.Error(w, "Error fetching chats: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, "Error formating answer", http.StatusInternalServerError)
	}

}

// GetChatUsersHandler возвращает список юзеров для чата
func GetChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get chats comand get")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	chatIDStr := r.URL.Query().Get("chat_id")
	if chatIDStr == "" {
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid chat_id", http.StatusBadRequest)
		return
	}
	chats, err := Repo.GetChatUsers(chatID)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chats); err != nil {
		http.Error(w, "Error formating answer", http.StatusInternalServerError)
	}

}

func UpdateLastReadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get update unread comand get")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		UserID int `json:"user_id"`
		ChatID int `json:"chat_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := Repo.UpdateLastRead(req.UserID, req.ChatID); err != nil {
		http.Error(w, "update error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUnreadCountsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get unread comand get")
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	counts, err := Repo.GetUnreadCount(userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}
