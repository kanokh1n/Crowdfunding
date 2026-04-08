# Crowdfunding Platform

REST API на Go + фронтенд (отдельный разработчик).

---

## Стек

| Слой | Технология |
|---|---|
| Backend | Go 1.25, Gin |
| ORM | GORM + pgx |
| База данных | PostgreSQL 16 |
| Кэш / токены | Redis 7 |
| Аутентификация | JWT (access 15 мин + refresh 7 дней) |
| Фронтенд | nginx (статика) |

---

## Быстрый старт (Docker)

### Требования
- [Docker](https://docs.docker.com/get-docker/) + [Docker Compose](https://docs.docker.com/compose/)

### Запуск

```bash
# Клонировать репозиторий
git clone <repo-url>
cd Crowdfunding

# Собрать и поднять все сервисы
make build
```

Сервисы:

| Сервис | URL |
|---|---|
| Backend API | http://localhost:8080 |
| Frontend | http://localhost:3000 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

### Основные команды

```bash
make up          # запустить в фоне
make down        # остановить
make logs        # смотреть логи всех сервисов
make logs-back   # логи только бэкенда
make restart     # перезапустить
make db-shell    # открыть psql
make redis-shell # открыть redis-cli
make help        # показать все команды
```

---

## Локальная разработка (без Docker)

### Требования
- Go 1.25+
- PostgreSQL 16
- Redis 7

### Настройка

```bash
cd back
cp .env.example .env
# Отредактируй .env под свои параметры
```

```bash
# Установить зависимости
make tidy

# Запустить сервер
make run
```

При старте GORM автоматически создаёт все таблицы.

---

## Переменные окружения (`back/.env`)

| Переменная | По умолчанию | Описание |
|---|---|---|
| `DB_URL` | `postgres://postgres:postgres@localhost:5432/crowdfunding?sslmode=disable` | Строка подключения к Postgres |
| `REDIS_URL` | `redis://localhost:6379` | Адрес Redis |
| `JWT_SECRET` | `dev-secret-change-in-production` | Секрет для подписи JWT (**обязательно сменить в проде**) |
| `PORT` | `8080` | Порт сервера |

---

## API документация

Swagger-спецификация: [`docs/swagger.yaml`](docs/swagger.yaml)

Открыть интерактивную документацию можно через [Swagger Editor](https://editor.swagger.io/) — импортируй файл `docs/swagger.yaml`.

Или локально через Docker:

```bash
docker run -p 8081:8080 \
  -e SWAGGER_JSON=/docs/swagger.yaml \
  -v $(pwd)/docs:/docs \
  swaggerapi/swagger-ui
```

Затем открой http://localhost:8081

---

## Структура проекта

```
Crowdfunding/
├── back/
│   ├── cmd/
│   │   └── main.go              # Точка входа, роутер
│   ├── internal/
│   │   ├── config/              # Загрузка переменных окружения
│   │   ├── database/            # Подключение Postgres и Redis
│   │   ├── model/               # Все модели БД
│   │   ├── middleware/          # JWT аутентификация, проверка роли
│   │   ├── moderation/          # AI-заглушка (замени на реальную нейронку)
│   │   └── handler/             # HTTP хендлеры по доменам
│   ├── .env.example
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── front/
│   ├── index.html               # Placeholder (заполняется фронтенд-разработчиком)
│   └── Dockerfile
├── docs/
│   └── swagger.yaml             # OpenAPI 3.0 спецификация
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Схема базы данных

```
users
  id, email, password, fio, description, profile_img, phone
  role (user|admin), is_verified, created_at, updated_at

email_tokens
  id, user_id → users, token, expires_at

categories
  id, title

projects
  id, user_id → users, title, description
  goal_amount, current_amount, project_img
  status (pending_ai|rejected_ai|pending_human|active|rejected|completed|cancelled)
  end_date, created_at, updated_at, deleted_at  ← soft delete

project_moderation  (1:1 с projects)
  id, project_id → projects
  ai_status (pending|passed|failed), ai_score (0.0–1.0), ai_flags (JSON), ai_checked_at
  human_status (pending|approved|rejected), moderator_id → users
  moderator_note, human_moderated_at

project_categories  (m2m)
  project_id → projects, category_id → categories

pledges
  id, user_id → users, project_id → projects
  amount, created_at

comments
  id, user_id → users, project_id → projects
  content, created_at

likes
  id, user_id → users, project_id → projects   ← unique(user_id, project_id)
  created_at

messages
  id, sender_id → users, recipient_id → users
  project_id → projects, title, content
  is_read, created_at
```

---

## Аутентификация (для фронтенда)

1. `POST /api/auth/register` — регистрация
2. `POST /api/auth/login` — получить `access_token` и `refresh_token`
3. Передавать `access_token` в заголовке: `Authorization: Bearer <token>`
4. Токен истекает через **15 минут** — вызвать `POST /api/auth/refresh` с `refresh_token`

```js
// Пример fetch
const res = await fetch('http://localhost:8080/api/projects', {
  headers: {
    'Authorization': `Bearer ${accessToken}`,
    'Content-Type': 'application/json',
  }
})
```

---

## Роли

| Роль | Возможности |
|---|---|
| `user` | Создавать проекты, делать pledges, комментировать, лайкать, писать сообщения |
| `admin` | Всё выше + управление пользователями, проектами и модерация через `/api/admin/*` |

---

## Модерация проектов

Каждый новый проект проходит два этапа перед публикацией.

### Этап 1 — AI-проверка (автоматически при создании)

Запускается синхронно сразу после `POST /api/projects`. Проверяет заголовок и описание на запрещённые ключевые слова и выставляет оценку (`ai_score` от 0 до 1).

| Результат | Статус проекта | Что происходит |
|---|---|---|
| Чисто (score ≥ 0.40, нет флагов) | `pending_human` | Проект уходит на проверку модератору |
| Нарушение (score < 0.40 или есть флаги) | `rejected_ai` | Проект отклонён автоматически |

Возможные `ai_flags`: `spam`, `gambling`, `drugs`, `weapons`, `adult_content`, `fraud`.

> Заглушка заменяется на реальную нейронку в `internal/moderation/ai.go` — функция `RunAICheck`.

### Этап 2 — Проверка модератором

Модератор работает через `/api/admin/moderation`:

```
GET  /api/admin/moderation                    # очередь проектов pending_human
GET  /api/admin/moderation/:project_id        # детали + AI результат
PATCH /api/admin/moderation/:project_id       # решение: approve / reject
POST /api/admin/moderation/:project_id/recheck  # повторный AI-прогон
```

Тело запроса для принятия решения:
```json
{
  "decision": "approve",
  "moderator_note": "Всё в порядке"
}
```
```json
{
  "decision": "reject",
  "moderator_note": "Проект содержит признаки рекламы запрещённых товаров"
}
```

### Жизненный цикл проекта

```
Создание → pending_ai
               │
         AI-проверка
               ├── не прошёл ──→ rejected_ai
               └── прошёл ────→ pending_human
                                     │
                               Модератор
                                     ├── reject ──→ rejected
                                     └── approve ─→ active
                                                      │
                                              (работающий проект)
                                                      ├── completed
                                                      └── cancelled
```

Публичный список проектов (`GET /api/projects`) возвращает **только `active`** проекты. Остальные статусы видны только через admin-эндпоинты.
