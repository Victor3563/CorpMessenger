// Заебашили работу с чатами
package repo

import (
	"errors"
)

type Conversation struct {
	ID        int    `json:"id"`
	Type      string `json:"type"` // "private" или "group"
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	CreatorID int    `json:"creator_id"`
}

type ConversationMember struct {
	ConversationID int    `json:"conversation_id"`
	UserID         int    `json:"user_id"`
	Role           string `json:"role"`
	JoinedAt       string `json:"joined_at"`
}

type Returns_chat_info struct {
	ConversationName string `json:"conversation_name"`
	Role             string `json:"role"`
}

// Работа с чатами
func (r *Repository) CreateConversation(convType, name string) (Conversation, error) {
	var conv Conversation
	query := `INSERT INTO conversations (type, name) VALUES ($1, $2) RETURNING id, type, name, created_at`
	err := r.DB.QueryRow(query, convType, name).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.CreatedAt)
	if err != nil {
		return conv, err
	}
	return conv, nil
}

func (r *Repository) DeleteConversation(convID int) error {
	query := `DELETE FROM conversations WHERE id = $1`
	res, err := r.DB.Exec(query, convID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// возвращает чаты пользователя
func (r *Repository) GetUserChats(userID int) ([]Conversation, error) {
	query := `SELECT c.id, c."type", c.name, c.created_at
		FROM conversations c
		JOIN conversation_members cm ON c.id = cm.conversation_id
		WHERE cm.user_id = $1
		ORDER BY c.created_at DESC`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //Гарантируем закрытие
	var chats []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.Type, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil

}

// возвращает пользователей чата
func (r *Repository) GetChatUsers(chatID int) ([]UserInfo, error) {
	query := `SELECT u.id, u.username, u.email, cm.role
		FROM users u
		JOIN conversation_members cm ON u.id = cm.user_id
		WHERE cm.conversation_id = $1
		ORDER BY u.username DESC`
	rows, err := r.DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //Гарантируем закрытие
	var chats []UserInfo
	for rows.Next() {
		var user UserInfo
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role); err != nil {
			return nil, err
		}
		chats = append(chats, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil

}

// Установка зависимостей между чатами и юзерами
func (r *Repository) AddMemberToConversation(convID, userID int, role string) error {
	query := `INSERT INTO conversation_members (conversation_id, user_id, role) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, convID, userID, role)
	return err
}

func (r *Repository) RemoveMemberFromConversation(convID, userID int) error {
	query := `DELETE FROM conversation_members WHERE conversation_id = $1 AND user_id = $2`
	res, err := r.DB.Exec(query, convID, userID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("member not found in conversation")
	}
	return nil
}

// возвращает количество непрочитанных сообщений из чата
func (r *Repository) GetUnreadCount(userID int) (map[int]int, error) {
	query := `
		SELECT m.chat_id, COUNT(*) AS unread_count
		FROM messages m
		JOIN conversation_members cm ON m.chat_id = cm.conversation_id
		LEFT JOIN unread_messages lr ON lr.user_id = $1 AND lr.chat_id = m.chat_id
		WHERE cm.user_id = $1
		  AND (lr.last_read_at IS NULL OR m.created_at > lr.last_read_at)
		  AND m.sender_id != $1
		  AND m.deleted = false
		GROUP BY m.chat_id
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var chatID, count int
		if err := rows.Scan(&chatID, &count); err != nil {
			return nil, err
		}
		result[chatID] = count
	}
	return result, nil
}

func (r *Repository) UpdateLastRead(userID, chatID int) error {
	query := `INSERT INTO unread_messages (user_id, chat_id, last_read_at)
			  VALUES ($1, $2, NOW())
			  ON CONFLICT (user_id, chat_id) DO UPDATE SET last_read_at = NOW()`
	_, err := r.DB.Exec(query, userID, chatID)
	return err
}
