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


Работа с сервером:
    Подключение к таблице:
    psql -h localhost -U user -d messenger

    Просмотр таблицы:
    SELECT * FROM messages;

    Команды для юзеров:
        Создание:
            curl -i -X POST -H "Content-Type: application/json" \
            -d '{"username": "testuser1", "password": "secret", "email": "user1@example.com"}' \
            http://localhost:8080/register

        Auth:
            curl -i -X POST -H "Content-Type: ation/json" -d '{"username" : "testuser1", "password": "secret"}' http://localhost:8080/auth

        Update:
            curl -i -X PUT -H "Content-Type: ation/json" -d '{"id": 5,"username" : "new_testuser1", "password": "new_secret","email": "new_user1@example.com"}' http://localhost:8080/updateUser


        Удаление:
            curl -i -X DELETE -H "Content-Type: application/json" \
            -d '{"id": !!!ID юзера!!!}' \
            http://localhost:8080/deleteUser

    Команды для чатов:
        Создание:
            curl -i -X POST -H "Content-Type: application/json" \
            -d '{"type": "private", "name": ""}' \
            http://localhost:8080/createChat

        Удаление:
            curl -i -X DELETE -H "Content-Type: application/json" \
            -d '{"id": !!!ID_чата!!!}' \
            http://localhost:8080/deleteChat

        Добавление в чат:
            curl -i -X POST -H "Content-Type: application/json" -d '{"conversation_id": 4, "user_id": 5, "role": "member"}' http://localhost:8080/addMember

        Удаление из чата:
            curl -i -X DELETE -H "Content-Type: application/json" \
            -d '{"conversation_id": !!!ID_чата!!!, "user_id": !!!ID_пользователя!!!}' \
            http://localhost:8080/removeMember

    Команды для сообщений:
        Подключение пользователя к чату    
            wscat -c "ws://localhost:8080/ws?chat_id=1&user_id=1"
        после этого отправляем сообщение
            {"content": "Сообщение"}


Пояснение для Музы:
Пока план такой, если пользователь активен, мы вкидываем в websocket только те чаты, в которых пользователь проявил активность. Соответственно в них сообщения обмениваются мгновенно. Остальные остаются в фоне.Если полбзователь офлайн, он автоматом в фоне. В фоне пользоваетлю показывается что в чате есть новое сообщение (и например можно еще количество), (предположительно для этого нужна новая табличка), если пользователь открывает его мы прогружаем сообщения.Как только пользователь выходит из фона мы подгружаем только кол-во новых сообщений в чатах, остальное подгружаем по запросу пользователя(когда он открывает непосредственно чат). Как работает кэш. Мы отправляем пользователю последние 20 сообщений в 10 последних чатах в которых он сидел. Как только пользователь подключается к чату, мы прогружаем кэш, если он есть (чтобы даже офлайн были сообщения), после чего через бд проверяем его валидность и обновляем в случае чего. Разделение на групповой/личный чат.Есть смысл забить на реализацию этого на сервене и просто убрать кнопку добавления новых участников в самом приложении в случаи если чат личный. Пока я так представляю структуру.

Сейчас у меня есть желание начать писать вэб, тк очень сильно не хватает возможности норм тестирования (У меня  оказывается удаление кринге работало, а я этого не понял, хз как), конечно прикольно в 3 терминалах тыкаться, но вэб будто сильно поможет. Как минимум у меня все готово для написания auth,  и там уже паралельно можно дописывать, смотря на то, чего не хватает