package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMessagesHandler_MissingChatID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/getMessage", nil)
	w := httptest.NewRecorder()
	GetMessagesHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "Необходим параметр chat_id\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}

func TestDeleteMessageHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/deleteMessage", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	DeleteMessageHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "неверный JSON\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}

func TestUpdateUserHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/updateUser", nil)
	w := httptest.NewRecorder()
	UpdateUserHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d; got %d", http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestDeleteUserHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/deleteUser", nil)
	w := httptest.NewRecorder()
	DeleteUserHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d; got %d", http.StatusMethodNotAllowed, res.StatusCode)
	}
}

func TestFindUserHandler_MissingUsername(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/findUser", nil)
	w := httptest.NewRecorder()
	FindUserHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "user_id is required\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}

func TestFindUserbyIDHandler_MissingUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/findUserbyID", nil)
	w := httptest.NewRecorder()
	FindUserbyIDHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "user_id is required\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}

func TestGetUserChatsHandler_MissingUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/getChats", nil)
	w := httptest.NewRecorder()
	GetUserChatsHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "user_id is required\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}

func TestGetUnreadCountsHandler_InvalidUserID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unreadCounts?user_id=abc", nil)
	w := httptest.NewRecorder()
	GetUnreadCountsHandler(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
	body := w.Body.String()
	expected := "invalid user_id\n"
	if body != expected {
		t.Errorf("expected body %q; got %q", expected, body)
	}
}
