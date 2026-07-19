# 🏦 OnlineBank API

<p align="left">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Echo-00C7B7?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white" />
  <img src="https://img.shields.io/badge/JWT-000000?style=flat-square&logo=jsonwebtokens&logoColor=white" />
</p>

REST API онлайн-банка на **Go** и фреймворке **Echo**. Реализованы регистрация и вход пользователей,
JWT-авторизация, проверка баланса и переводы средств между счетами. Данные хранятся в **PostgreSQL**,
миграции применяются через `golang-migrate`. Проект построен по принципам чистой архитектуры.

---

## 🧱 Архитектура

```
cmd/                     — точка входа
internal/
├── config/              — чтение конфигурации из переменных окружения
├── db/                  — подключение к PostgreSQL
├── handlers/            — HTTP-обработчики (Echo)
├── middlewarecustom/    — middleware проверки JWT / баланса
├── models/              — модели user и transaction
├── repository/          — слой доступа к данным
├── servises/
│   ├── jwtservice/      — генерация и проверка JWT
│   └── userservice/     — бизнес-логика пользователей и переводов
└── utils/
    ├── bind/            — привязка тела запроса
    └── hashpassword/    — хэширование паролей (bcrypt)
migrations/              — SQL-миграции (golang-migrate)
web/                     — статические страницы (register, dashboard)
```

---

## ⚙️ Технологии

- **Go** + [Echo](https://echo.labstack.com/) — веб-фреймворк
- **PostgreSQL** (`lib/pq`) — хранение пользователей и балансов
- **JWT** (`golang-jwt/jwt/v5`) — авторизация
- **bcrypt** (`golang.org/x/crypto`) — хэширование паролей
- **golang-migrate** — миграции БД
- **testify** — тесты

---

## 📚 API-эндпоинты

| Метод | URL | Описание | Авторизация |
|-------|-----|----------|:-----------:|
| POST  | `/api/register` | Регистрация нового пользователя | — |
| POST  | `/api/login`    | Вход, установка JWT в cookie     | — |
| GET   | `/api/balance`  | Проверка текущего баланса        | ✅ JWT |
| POST  | `/api/transfer` | Перевод средств другому пользователю по email | ✅ JWT |

> При регистрации каждому пользователю начисляется стартовый баланс **1000**.

---

## 🚀 Запуск

### 1. Требования

- Go 1.24+
- PostgreSQL
- [`golang-migrate`](https://github.com/golang-migrate/migrate) (для миграций)

### 2. Переменные окружения

| Переменная | Описание | Пример |
|---|---|---|
| `PORT` | Порт сервера (по умолчанию `8080`) | `8080` |
| `DATABASE_URL` | Строка подключения к PostgreSQL | `postgres://postgres:1234@localhost:5432/postgres?sslmode=disable` |
| `SECRET_KEY_JWT` | Секретный ключ для подписи JWT | `your-secret-key` |

### 3. Миграции и запуск

```bash
# применить миграции
make migrate-up

# запустить сервер
make run
```

> В `makefile` уже прописаны команды `run`, `migrate-up` и `migrate-down`.
> Значения строки подключения и секретного ключа задавайте под своё окружение.

---

## 🔎 Примеры запросов

```bash
# Регистрация
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Ivan","last_name":"Petrov","phone":"79990000000","email":"ivan@mail.com","password":"12345678"}'

# Вход (JWT возвращается в cookie)
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -c cookie.txt \
  -d '{"email":"ivan@mail.com","password":"12345678"}'

# Проверка баланса
curl http://localhost:8080/api/balance -b cookie.txt

# Перевод средств
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" -b cookie.txt \
  -d '{"to":"receiver@mail.com","amount":100}'
```

---

## 🧪 Тесты

```bash
go test ./...
```

---

<p align="center"><sub>Сделано с ❤️ и Go · <a href="https://github.com/GoCoreDevelopment">GoCoreDevelopment</a></sub></p>
