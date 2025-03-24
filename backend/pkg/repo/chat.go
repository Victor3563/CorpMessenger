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
