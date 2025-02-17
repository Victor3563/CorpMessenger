### Corp Messenger
Запускаем миграцию:
migrate -path internal/migrations -database "postgres://user:password@localhost:5432/messenger?sslmode=disable" up
-path internal/migrations – путь к файлам миграций.
-database "postgres://..." – строка подключения к базе.
up – применить миграции.
Если все прошло успешно, можно проверить, есть ли таблица:
docker exec -it messenger_db psql -U user -d messenger -c "\dt"
Если нужно отменить последнюю миграцию:
migrate -path internal/migrations -database "postgres://user:password@localhost:5432/messenger?sslmode=disable" down 1

Запуск докера:
docker-compose up -d
Проверить статус:
docker ps

Подключение к бд:
docker exec -it messenger_db psql -U user -d messenger

Отключить докер:
docker-compose down
