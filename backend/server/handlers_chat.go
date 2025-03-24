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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conv)
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
