# Tinker Backend (Go)

Go backend для социальной сети Tinker.

## Структура проекта

```
backend/
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── config/         # Конфигурация
│   ├── database/       # Модели БД и подключение
│   ├── handlers/       # HTTP handlers
│   ├── services/       # Бизнес-логика
│   ├── middleware/     # Middleware (auth, logging, CORS)
│   ├── dto/           # Структуры запросов/ответов
│   └── utils/         # Утилиты
├── migrations/         # SQL миграции
├── seed/              # Seed данные
└── Dockerfile         # Docker конфигурация
```

## Запуск через Docker Compose

```bash
docker-compose up backend2
```

## Локальная разработка

### Требования

- Go 1.21+
- PostgreSQL 15+
- golang-migrate (для миграций)

### Установка зависимостей

```bash
go mod download
```

### Настройка переменных окружения

Создайте `.env` файл или установите переменные:

```bash
export DATABASE_URL="postgresql://postgres:postgres@localhost:5444/tnews"
export JWT_SECRET="your-secret-key"
export PORT="3001"
export LOG_LEVEL="info"
export SEED_DB="false"
```

### Применение миграций

```bash
make migrate
```

или вручную:

```bash
migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5444/tnews?sslmode=disable" up
```

### Запуск seed данных

```bash
make seed
```

или:

```bash
go run ./seed/seed.go
```

### Запуск сервера

```bash
make run
```

или:

```bash
go run ./cmd/server/main.go
```

### Сборка

```bash
make build
```

## API Endpoints

Все эндпоинты доступны по префиксу `/api`:

- `POST /api/auth/login` - Аутентификация
- `GET /api/users` - Список пользователей
- `POST /api/users` - Создание пользователя
- `GET /api/users/:userId` - Получение пользователя
- `PATCH /api/users/:userId` - Обновление профиля (требует auth)
- `GET /api/users/:userId/posts` - Посты пользователя
- `POST /api/users/:userId/posts` - Создание поста (требует auth)
- `POST /api/users/:userId/follow` - Подписка (требует auth)
- `DELETE /api/users/:userId/follow` - Отписка (требует auth)
- `GET /api/users/:userId/following` - Список подписок
- `GET /api/posts` - Все посты
- `DELETE /api/posts/:postId` - Удаление поста (требует auth)
- `POST /api/posts/:postId/likes` - Лайк (требует auth)
- `DELETE /api/posts/:postId/likes` - Анлайк (требует auth)
- `GET /api/posts/:postId/comments` - Комментарии к посту
- `POST /api/posts/:postId/comments` - Создание комментария (требует auth)
- `DELETE /api/comments/:commentId` - Удаление комментария (требует auth)
- `GET /api/feed` - Лента постов (требует auth)
- `GET /api/search?query=...&type=users|posts` - Поиск

## Аутентификация

Используется JWT токены. После успешного логина токен передается в заголовке:

```
Authorization: Bearer <token>
```

## Технологии

- **Gin** - веб-фреймворк
- **GORM** - ORM
- **golang-jwt/jwt** - JWT аутентификация
- **go-playground/validator** - валидация
- **slog** - логирование (встроенное в Go 1.21+)
- **golang-migrate** - миграции БД
- **bcrypt** - хеширование паролей

