// Это черновой вариант,реализовано с излишеком/недостатоком/дублированием подробнее смотри в websocket
package repo

import (
	"errors"
	"fmt"
	"time"
)

// WSMessage структура сообщения

type WSMessage struct {
	ChatID   int    `json:"chat_id"`
	SenderID int    `json:"sender_id"`
	Content  string `json:"content"`
}

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	SenderID  int       `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"deleted"`
}

// Сохраняет новое сообщение в базе данных и возвращает его.
func (r *Repository) AddMessage(chatID, senderID int, content string) (Message, error) {
	var msg Message
	query := `INSERT INTO messages (chat_id, sender_id, content, created_at, deleted)
		VALUES ($1, $2, $3, NOW(), false)
		RETURNING id, chat_id, sender_id, content, created_at, deleted;`
	err := r.DB.QueryRow(query, chatID, senderID, content).
		Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.Deleted)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

// Сохраняет много сообщений в базе данных.
func (r *Repository) BatchInsertMessages(messages []WSMessage) error {
	if len(messages) == 0 {
		return nil
	}
	query := "INSERT INTO messages (chat_id, sender_id, content, created_at, deleted) VALUES "
	args := []interface{}{}
	for i, msg := range messages {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, NOW(), false)", i*3+1, i*3+2, i*3+3)
		args = append(args, msg.ChatID, msg.SenderID, msg.Content)
	}
	_, err := r.DB.Exec(query, args...)
	return err
}

// возвращает последние N сообщений для указанного чата, упорядоченные по возрастанию времени.
func (r *Repository) GetMessages(chatID, limit int) ([]Message, error) {
	// Заебался пистаь запросы
	query := `SELECT id, chat_id, sender_id, content, created_at, deleted
		FROM messages
		WHERE chat_id = $1 AND deleted = false
		ORDER BY created_at ASC
		LIMIT $2;`
	rows, err := r.DB.Query(query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //Гарантируем закрытие
	messages := make([]Message, 0, limit)
	for rows.Next() {
		msg := Message{}
		err = rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.Deleted)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	//Переворачиваем срез, чтобы вернуть сообщения в порядке возрастания времени. Пока считаю это костылем. Надо прогуглить!!!
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

// Выполняет soft delete сообщения (устанавливает deleted = true).
// Удаление может выполнить только отправитель сообщения!!!
func (r *Repository) DeleteMessage(messageID, senderID int) error {
	//Уместил все в один запрос
	query := `UPDATE messages 
			SET deleted = true
			WHERE id = $1
				AND deleted = false 
            	AND sender_id = $2`
	res, err := r.DB.Exec(query, messageID, senderID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Message not found or haven't ability to delete it")
	}
	return nil
}
