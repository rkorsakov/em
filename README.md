# Subscription Service

REST API сервис для управления онлайн-подписками пользователей.

## Запуск

Для запуска сервиса используйте следующую команду:

```bash
  make build
```

Сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## API Endpoints

### Подписки

| Метод   | Эндпоинт                              | Описание                     |
|---------|---------------------------------------|------------------------------|
| `POST`  | `/api/v1/subscriptions`              | Создать подписку             |
| `GET`   | `/api/v1/subscriptions`              | Получить список подписок     |
| `GET`   | `/api/v1/subscriptions/{id}`         | Получить подписку по ID      |
| `PUT`   | `/api/v1/subscriptions/{id}`         | Обновить подписку            |
| `DELETE`| `/api/v1/subscriptions/{id}`         | Удалить подписку             |

### Агрегация

| Метод   | Эндпоинт                              | Описание                                    |
|---------|---------------------------------------|---------------------------------------------|
| `GET`   | `/api/v1/summary`                    | Получить суммарную стоимость подписок за период |

### Утилиты

| Метод   | Эндпоинт                              | Описание                     |
|---------|---------------------------------------|------------------------------|
| `GET`   | `/health`                            | Health check                 |
| `GET`   | `/swagger/index.html`                | Swagger документация         |

## Примеры запросов

### Создание подписки

```bash
  curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "15-01-2025"
  }'
```

### Получение суммы за период

```bash
  curl  "http://localhost:8080/api/v1/summary?start_period=01-01-2025&end_period=31-12-2025"
```

### Фильтрация по пользователю

```bash
  curl  "http://localhost:8080/api/v1/subscriptions?user_id=user-1"
```