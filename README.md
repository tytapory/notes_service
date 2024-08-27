# Заметки REST API

Это простой сервер позволяющий пользователям регистрироваться, создавать заметки и получать их список.

## API Endpoints

### 1. Регистрация пользователя (`/register`)

**Описание:** Регистрация нового пользователя.

**Запрос:**
```bash
curl -X POST http://localhost:8080/register \
-H "Content-Type: application/json" \
-d '{
  "credentials": {
    "username": "user1",
    "password": "password123"
  }
}'
```

**Ответ (успех):**
```json
{
  "message": "User successfully registered"
}
```

**Ответ (ошибка):**
```json
{
  "error": "Error message describing the issue"
}
```

### 2. Добавление заметки (`/push_note`)

**Описание:** Добавление новой заметки для зарегистрированного пользователя.

**Запрос:**
```bash
curl -X POST http://localhost:8080/push_note \
-H "Content-Type: application/json" \
-d '{
  "credentials": {
    "username": "user1",
    "password": "password123"
  },
  "note": "This is a new note"
}'
```

**Ответ (успех):**
```json
{
  "message": "Note successfully created"
}
```

**Ответ (ошибка):**
```json
{
  "error": "Error message describing the issue"
}
```

### 3. Получение заметок (`/get_notes`)

**Описание:** Получение всех заметок для зарегистрированного пользователя.

**Запрос:**
```bash
curl -X POST http://localhost:8080/get_notes \
-H "Content-Type: application/json" \
-d '{
  "credentials": {
    "username": "user1",
    "password": "password123"
  }
}'
```

**Ответ (успех):**
```json
{
  "notes": [
    "This is a new note",
    "This is another note"
  ]
}
```

**Ответ (ошибка):**
```json
{
  "error": "Error message describing the issue"
}
```

