package root

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
)

// парсим запрос и проверяем его адекватность (просто упростили код heandlers) челики
func ParserandValidByName(r *http.Request) (repo.User, error) {
	var req repo.User

	if r.Method != http.MethodPost {
		return req, errors.New("invalid request method by user")
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON in user")
	}

	if req.Name == "" || req.Password == "" {
		return req, errors.New("missing required fields")
	}

	return req, nil
}

func ParserandValidByID(r *http.Request) (repo.User, error) {
	var req repo.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON in user")
	}

	if req.ID == 0 {
		return req, errors.New("missing required ID in user")
	}

	return req, nil
}

// Чатики
func ParserConversationAdd(r *http.Request) (repo.Conversation, error) {
	var req repo.Conversation
	if r.Method != http.MethodPost {
		return req, errors.New("invalid request method conversation")
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON in conversation")
	}
	if req.Type == "" || req.Name == "" {
		return req, errors.New("missing required fields")
	}
	return req, nil
}

func ParserConversationDelete(r *http.Request) (repo.Conversation, error) {
	var req repo.Conversation
	if r.Method != http.MethodDelete {
		return req, errors.New("invalid request method conversation")
	}
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		return req, errors.New("missing required id")
	}
	var err error
	req.ID, err = strconv.Atoi(idStr)
	if err != nil {
		return req, errors.New("invalid conversation id")
	}
	return req, nil

}

func ParserAddToConversation(r *http.Request) (repo.ConversationMember, error) {
	var req repo.ConversationMember
	if r.Method != http.MethodPost {
		return req, errors.New("invalid request method conversation")
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON in conversation")
	}
	if req.Role == "" {
		req.Role = "member"
	}
	return req, nil
}

func ParserRemoveFromConversation(r *http.Request) (repo.ConversationMember, error) {
	var req repo.ConversationMember
	if r.Method != http.MethodDelete {
		return req, errors.New("invalid request method conversation")
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON in conversation")
	}
	return req, nil
}
