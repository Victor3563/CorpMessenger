package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Victor3563/CorpMessenger/config"
	"github.com/Victor3563/CorpMessenger/pkg/repo"
	"github.com/Victor3563/CorpMessenger/server"
	_ "github.com/lib/pq"
)

// corsMiddleware оборачивает обработчик и добавляет необходимые заголовки CORS.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любого источника. В продакшене можно указать конкретный URL, например, "http://localhost:5173"
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Разрешенные HTTP-методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Разрешенные заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это предзапрос с методом OPTIONS, возвращаем статус OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func start_rep() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port,
		config.Database.DBName, config.Database.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализация репозитория и передача его в глобальную переменную обработчиков
	repoInstance := repo.NewRepository(db)
	server.InitRoutes(repoInstance)

	addr := fmt.Sprintf(":%d", config.Server.Port)
	fmt.Printf("Server started at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, corsMiddleware(http.DefaultServeMux)))
}

func main() {
	start_rep()
}
