# AvitoTech-Task
Rest API для сервиса управления баннерами.

## Быстрый старт
```bash
make
```
```bash
make migrate-up
```
```bash
make testdata
```
*[Swagger](http://localhost:8080/swagger)*

* Токен пользователя: **user**
* Токен админа: **admin**

Остановить приложение:
```bash
make stop
```

## Описание Makefile
1. Собирает и запускает приложение
```bash
make all
```

2. Собирает Docker-образ приложения
```bash
make build
```

3. Запускает Docker-контейнер с приложением
```bash
make run
```

4. Останавливает Docker-контейнер с приложением
```bash
make stop
```

5. Запускает Docker-контейнер с базой данных
```bash
make db-start
```

6. Останавливает Docker-контейнер с базой данных
```bash
make db-stop
```

7. Применяет миграции базы данных вверх
```bash
make migrate-up
```

8. Откатывает миграции базы данных
```bash
make migrate-down
```

9. Загружает тестовые данные в базу данных
```bash
make testdata
```

10. Запускает интеграционные тесты
```bash
make integration-test
```
