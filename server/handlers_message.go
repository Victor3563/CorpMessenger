// Черновик, нужен рефакторинг
package server

import (
	"encoding/json"
	"net/http"
)

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
