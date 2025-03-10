// Это черновой вариант,реализовано с излишеком/недостатоком/дублированием подробнее смотри в websocket
package repo

import (
	"errors"
	"time"
)

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

// возвращает последние N сообщений для указанного чата, упорядоченные по возрастанию времени.
func (r *Repository) GetMessages(chatID, limit int) ([]Message, error) {
	// Заебался пистаь запросы
	query := `SELECT id, chat_id, sender_id, content, created_at, deleted
		FROM messages
		WHERE chat_id = $1 AND deleted = false
		ORDER BY created_at DESC
		LIMIT $2;`
	rows, err := r.DB.Query(query, chatID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //Гарантируем закрытие

	messages := []Message{}
	for rows.Next() {
		var msg Message
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
	var dbSenderID int
	// Я ебал эти запросы, забыл FROM, дебажил 30 минут
	querySelect := `SELECT sender_id FROM messages WHERE id = $1`
	err := r.DB.QueryRow(querySelect, messageID).Scan(&dbSenderID)
	if err != nil {
		return err
	}
	if dbSenderID != senderID {
		return errors.New("нельзя удалить сообщение: неверный отправитель")
	}

	query := `UPDATE messages SET deleted = true WHERE id = $1`
	res, err := r.DB.Exec(query, messageID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("сообщение не найдено или уже удалено")
	}
	return nil
}
