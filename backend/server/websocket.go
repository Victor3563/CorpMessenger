//Пора перекинуть все принты в логирование буду заниматься постепенно, по мере рефакторинга
//Кэширование будет проходить на стороне юзера, так что сервер просто возращает ему по запросу N сообщений.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Victor3563/CorpMessenger/pkg/repo"
	"github.com/gorilla/websocket"
)

// Client подключаем к WS
type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	ChatID    int
	UserID    int
	closeOnce sync.Once //Хотим закрываться только раз
}

// Закрывашка клиентов
func (c *Client) safeCloseSend() {
	c.closeOnce.Do(func() {
		close(c.Send)
	})
}

// Hub управляет всеми подключениями, распределяет входящие сообщения по чатам.
type Hub struct {
	//Я отказался от хранения клиентов в бд в пользу подобной структуры, так как обращение к бд ресурсно затратно
	Clients    map[int]map[*Client]bool //Пофакту set клиентов по чат id
	Broadcast  chan repo.WSMessage      // Канал для сообщений
	Register   chan *Client             // Каналы для работы с клиентами
	Remove     chan *Client             // Каналы для работы с клиентами
	BatchQueue chan repo.WSMessage

	Mutex sync.RWMutex //Используем RWMutex так как писать будем
	//  не так часто как читать, и нам хочется для скорости разделять эти процессы
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
		Clients:    make(map[int]map[*Client]bool),
		Broadcast:  make(chan repo.WSMessage),
		Register:   make(chan *Client),
		Remove:     make(chan *Client),
		BatchQueue: make(chan repo.WSMessage, 1000), // Буфер для сообщений батча
	}
}

// Может стоит пакетную вставку вообще вынести в другой файл. Пока оставлю тут, но я готов...
// Cобирает сообщения из BatchQueue и выполняет пакетную вставку в базу.
func (h *Hub) InsertProcessor() {
	// Таймер для таймаута (например, каждые 5 секунд (стоит ли вынести это в отдельный класс Settings?))
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()
	batch := make([]repo.WSMessage, 0, 100)
	for {
		select {
		case msg := <-h.BatchQueue:
			batch = append(batch, msg)
			// Если накопилось 100 сообщений, обрабатываем батч
			if len(batch) >= 100 {
				if err := Repo.BatchInsertMessages(batch); err != nil {
					fmt.Printf("Batch insert error: %v", err)
					return
				}
				batch = batch[:0]
			}
		case <-timer.C:
			// Если прошло время, а в пакете есть сообщения, обрабатываем их.
			if len(batch) > 0 {
				if err := Repo.BatchInsertMessages(batch); err != nil {
					fmt.Printf("Batch insert error: %v", err)
					return
				}
				batch = batch[:0]
			}
		}
	}
}

func (h *Hub) Run() {
	go h.InsertProcessor()
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			if _, status_ok := h.Clients[client.ChatID]; !status_ok {
				h.Clients[client.ChatID] = make(map[*Client]bool)
			}
			h.Clients[client.ChatID][client] = true
			h.Mutex.Unlock()
			fmt.Printf("Client connect: chat_id = %d, user_id = %d", client.ChatID, client.UserID)
		case client := <-h.Remove:
			h.Mutex.Lock()
			if clients, status_ok := h.Clients[client.ChatID]; status_ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					client.safeCloseSend()
					fmt.Printf("Отключен клиент: chat_id=%d, user_id=%d", client.ChatID, client.UserID)
					if len(clients) == 0 {
						delete(h.Clients, client.ChatID)
					}
				}
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			payload, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Ошибка в преобразовании сообщения в джсон: %v", err)
				continue
			}
			// Рассылаем сообщение всем подключенным клиентам в чате, кроме отправителя.
			// Вот тут (или не тут)) )нужно отправлять уведы не подключеным клиентам выставляя в новой табличке соответсвующие поля, написать отдельную функцию
			// которая еще пробежит и скажет всем что нью увед есть
			h.Mutex.RLock()
			if clients, ok := h.Clients[message.ChatID]; ok {
				for client := range clients {
					select {
					case client.Send <- payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.Mutex.RUnlock()

			//Сохраняем в бд только после отправки всем(Отправляем в поток для сохранения)
			h.BatchQueue <- message
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
		var msg repo.WSMessage
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
