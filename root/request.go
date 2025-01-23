package root

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
)

// парсим запрос и проверяем его адекватность (просто упростили код heandlers)
func ParserandValidByName(r *http.Request) (repo.User, error) {
	var req repo.User

	if r.Method != http.MethodPost {
		return req, errors.New("invalid request method")
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON")
	}

	if req.Name == "" || req.Password == "" {
		return req, errors.New("missing required fields")
	}

	return req, nil
}

func ParserandValidByID(r *http.Request) (repo.User, error) {
	var req repo.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, errors.New("invalid JSON")
	}

	if req.ID == "" {
		return req, errors.New("missing required ID")
	}

	return req, nil
}
