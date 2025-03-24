package server

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "недопустимый метод запроса", http.StatusMethodNotAllowed)
	}
	queryParams := r.URL.Query()
	chatIDParam := queryParams.Get("chat_id")
	if chatIDParam == "" {
		http.Error(w, "Необходим параметр chat_id", http.StatusBadRequest)
		return
	}

	chatID, err := strconv.Atoi(chatIDParam)
	if err != nil {
		http.Error(w, "Некорректный формат chat_id", http.StatusBadRequest)
		return
	}
	messages, err := Repo.GetMessages(chatID, 20)
	if err != nil {
		http.Error(w, "Ошибка при получении сообщений: ", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
	}
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "недопустимый метод запроса", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		MessageID int `json:"message_id"`
		SenderID  int `json:"sender_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный JSON", http.StatusBadRequest)
		return
	}
	if err := Repo.DeleteMessage(req.MessageID, req.SenderID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Сообщение успешно удалено"))
}
