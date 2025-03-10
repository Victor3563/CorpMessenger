//Это черновое представление файла, тут есть упрощения в пользу тестирования и дублирование структур,
// реализованных в других файлах, которое стоит пересмотреть при рефакторинге
//Пора перекинуть все принты в логирование буду заниматься постепенно, по мере рефакторинга

// Я ебал этот вебсокет, 4 дня безостановчной дрочки, потому что дохуя чего оказывается надо использовать, а инфа иногда не столь тривиально гуглится
// некоторые моменты все еще не верно работают.
// Базовый функционал тестировал перед отправкой(создание юзеров, чата, обмен сообщениями, удаление сообщений)

//Кэширование будет проходить на стороне юзера, так что сервер просто возращает ему по запросу N сообщений.
//Удаление пока работает через http, наверное так и останется, исходя из того что удалять юзеры будут реже чем писать. Надо обновлять сообщения, выводимые эзерам после удаления,
// пока это реализовано просто как флаг, но мы еще подумаем

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// WSMessage структура сообщения
type WSMessage struct {
	ChatID   int    `json:"chat_id"`
	SenderID int    `json:"sender_id"`
	Content  string `json:"content"`
}

// Client подключаем к WS
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	ChatID int
	UserID int
}

// Hub управляет всеми подключениями, распределяет входящие сообщения по чатам.
type Hub struct {
	//Я отказался от хранения клиентов в бд в пользу подобной структуры, так как обращение к бд ресурсно затратно
	Clients   map[int]map[*Client]bool //Пофакту set клиентов по чат id
	Broadcast chan WSMessage           // Канал для сообщений
	Register  chan *Client             // Каналы для работы с клиентами
	Remove    chan *Client             // Каналы для работы с клиентами
}

var wsUpgrader = websocket.Upgrader{
	//Допилить аутенфикацию и остальные настройки, пока оставлю пустыми для удобства тестирования
	// HandshakeTimeout:  Second, // Таймаут рукопожатия
	// ReadBufferSize:    1024,             // Размер буфера чтения
	// WriteBufferSize:   1024,             // Размер буфера записи
	CheckOrigin: func(r *http.Request) bool {
		// Здесь можно добавить проверку источника запроса, домен и валидность токена/пользователя если надо...
		return true
	},
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[int]map[*Client]bool),
		Broadcast: make(chan WSMessage),
		Register:  make(chan *Client),
		Remove:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, status_ok := h.Clients[client.ChatID]; !status_ok {
				h.Clients[client.ChatID] = make(map[*Client]bool)
			}
			h.Clients[client.ChatID][client] = true
			fmt.Printf("Client connect: chat_id = %d, user_id = %d", client.ChatID, client.UserID)
		case client := <-h.Remove:
			if clients, status_ok := h.Clients[client.ChatID]; status_ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					fmt.Printf("Отключен клиент: chat_id=%d, user_id=%d", client.ChatID, client.UserID)
					if len(clients) == 0 {
						delete(h.Clients, client.ChatID)
					}
				}
			}
		case message := <-h.Broadcast:
			//Сохраняем в бд
			savedMessage, err := Repo.AddMessage(message.ChatID, message.SenderID, message.Content)
			if err != nil {
				fmt.Printf("Ошибка сохранения сообщения: %v", err)
				continue
			}
			payload, err := json.Marshal(savedMessage)
			if err != nil {
				fmt.Printf("Ошибка в преобразовании сообщения в джсон: %v", err)
				continue
			}
			// Рассылаем сообщение всем подключенным клиентам в чате, кроме отправителя.
			if clients, ok := h.Clients[message.ChatID]; ok {
				for client := range clients {
					if client.UserID == message.SenderID {
						continue
					}
					select {
					case client.Send <- payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

// Читаем сообщения от клиента и отправляет их в Hub.
func (c *Client) ReadclMes() {
	//обработка отключений
	defer func() {
		c.Hub.Remove <- c
		c.Conn.Close()
	}()
	for {
		var msg WSMessage
		if err := c.Conn.ReadJSON(&msg); err != nil {
			fmt.Printf("Ошибка чтения JSON: %v", err)
			break
		}
		// Гарантируем, что сообщение относится к чату и отправителю, указанным при подключении.
		// Пока просто в тупую. Это вроде оптимальное решение, после обдумий, я не смог придумать,
		// почему так не стоит делать. Конечно в случае чего код перепишится, но пока оптимальная страта
		msg.ChatID = c.ChatID
		msg.SenderID = c.UserID
		c.Hub.Broadcast <- msg
	}
}

// vОтправляем сообщения клиенту.
func (c *Client) WriteclMes() {
	defer c.Conn.Close()
	for message := range c.Send {
		// Отправляем сообщение через WebSocket
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	// канал закрыт = отправляем сообщение о закрытии соединения
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
}

// Обрабатываем запросы на подключение по WebSocket.Ждем chat_id и user_id
func WSHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatIDStr := r.URL.Query().Get("chat_id")
		userIDStr := r.URL.Query().Get("user_id")
		if chatIDStr == "" || userIDStr == "" {
			http.Error(w, "нет chat_id иди user_id", http.StatusBadRequest)
			return
		}
		chatID, err := strconv.Atoi(chatIDStr)
		if err != nil {
			http.Error(w, "неверный chat_id", http.StatusBadRequest)
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "неверный user_id", http.StatusBadRequest)
			return
		}
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Ошибка апгрейда WebSocket: %v", err)
			return
		}
		client := &Client{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			ChatID: chatID,
			UserID: userID,
		}
		hub.Register <- client
		// Запускаем горутины для чтения и записи.
		go client.WriteclMes()
		client.ReadclMes()
	}
}
