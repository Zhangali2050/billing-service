# 💳 Billing Service

Сервис для выставления счетов, приёма платежей, хранения истории и интеграции с [AirbaPay](https://doc.airbapay.kz).

---

## 📦 Возможности

- Выставление счетов для пользователей с ролями `student` и `parent`
- Интеграция с AirbaPay:
  - Добавление и удаление карт
  - Проведение платежей
  - Возвраты
  - Вебхуки (оплата, сохранение карт)
- Хранение истории платежей
- REST API с защитой по API-ключу
- Встроенный интерфейс для тестирования (HTML)

---

## 🚀 Запуск в Docker

### 1. Убедитесь, что установлены:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### 2. Запуск

```bash
docker-compose up --build

---

## 🌐 Интерфейс

### Открой в браузере:

http://localhost:8080/

---

## 🔐 API-ключ

### Все запросы должны содержать заголовок:

X-Api-Key: sandbox_123

---

## 📂 Структура

.
├── cmd/                   # Точка входа (main.go)
├── internal/
│   ├── handler/           # HTTP обработчики
│   ├── service/           # Бизнес-логика
│   ├── repository/        # Работа с БД
│   ├── model/             # Общие структуры
│   └── airba/             # Клиент для AirbaPay
├── static/                # Интерфейс для тестирования
├── .env                   # Конфигурация (переменные)
├── docker-compose.yml
├── Dockerfile
└── go.mod / go.sum

---

## 🛠 Примеры API-запросов

### 1.Создание счёта

POST /invoice
Content-Type: application/json
X-Api-Key: sandbox_123

{
  "role": "student",
  "user_id": "uuid",
  "amount": 7900,
  "quantity": 1
}

### 2.Получение истории платежей

GET /payments?role=student&user_id=uuid
X-Api-Key: sandbox_123

---

## 📝 Авторизация в AirbaPay

### 1.В .env укажите:

AIRBA_USER=...
AIRBA_PASSWORD=...
AIRBA_TERMINAL_ID=...

### 2.Для вебхуков потребуется настроить публичный адрес (например, через ngrok).

---

## 🗄 Миграции БД

### Проект автоматически использует SQL-файл schema.sql. Убедитесь, что расширение pgcrypto подключено, если используете gen_random_uuid():

CREATE EXTENSION IF NOT EXISTS pgcrypto;

---

## 🧪 Тестовые данные

### AirbaPay предоставляет тестовые карты и сценарии.




Запустить докер
docker-compose up --build -d

Выключить докер
docker-compose down
