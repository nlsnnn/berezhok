# Бережок — REST API Contract

> **Версия API:** v1  
> **Base URL:** `https://api.berezhok.ru/v1`  
> **Формат:** JSON  
> **Аутентификация:** Bearer Token (JWT)

---

## Общие правила

### HTTP методы
- `GET` — получение данных
- `POST` — создание ресурса
- `PUT` — полное обновление ресурса
- `PATCH` — частичное обновление ресурса
- `DELETE` — удаление ресурса

### Стандартные HTTP статусы
- `200 OK` — успешный запрос
- `201 Created` — ресурс создан
- `204 No Content` — успешно, без тела ответа
- `400 Bad Request` — невалидные данные
- `401 Unauthorized` — не авторизован
- `403 Forbidden` — нет прав доступа
- `404 Not Found` — ресурс не найден
- `409 Conflict` — конфликт (например, дубликат)
- `422 Unprocessable Entity` — ошибка валидации
- `500 Internal Server Error` — ошибка сервера

### Формат ответа

**Успешный ответ:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Ответ с ошибкой:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Некорректные данные",
    "details": {
      "phone": "Неверный формат номера телефона"
    }
  }
}
```

### Пагинация

Для списков используется query параметры:
- `limit` — количество элементов (по умолчанию 20, макс 100)
- `offset` — смещение (по умолчанию 0)

**Пример запроса:**
```
GET /api/locations?limit=20&offset=40
```

**Пример ответа:**
```json
{
  "success": true,
  "data": {
    "items": [...],
    "pagination": {
      "total": 156,
      "limit": 20,
      "offset": 40,
      "has_more": true
    }
  }
}
```

---

## 1. Аутентификация

### 1.1 Отправка SMS кода (клиент)

**Endpoint:** `POST /auth/customer/send-code`

**Request:**
```json
{
  "phone": "+79001234567"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Код отправлен",
    "expires_in": 300
  }
}
```

**Errors:**
- `400` — некорректный формат телефона
- `429` — слишком много запросов (rate limit)

---

### 1.2 Вход по SMS коду (клиент)

**Endpoint:** `POST /auth/customer/login`

**Request:**
```json
{
  "phone": "+79001234567",
  "code": "123456"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "phone": "+79001234567",
      "name": ""
    }
  }
}
```

**Errors:**
- `401` — неверный код
- `400` — код истёк

---

### 1.3 Вход для партнёра

**Endpoint:** `POST /auth/partner/login`

**Request:**
```json
{
  "email": "owner@bakery.ru",
  "password": "SecurePassword123"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "650e8400-e29b-41d4-a716-446655440001",
      "email": "owner@bakery.ru",
      "name": "Иван Иванов",
      "role": "owner",
      "partner_id": "750e8400-e29b-41d4-a716-446655440002",
      "location_id": "850e8400-e29b-41d4-a716-446655440003"
    }
  }
}
```

**Errors:**
- `401` — неверные credentials
- `403` — аккаунт заблокирован

---

### 1.4 Вход для администратора

**Endpoint:** `POST /auth/admin/login`

**Request:**
```json
{
  "email": "admin@berezhok.ru",
  "password": "AdminPass123"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "950e8400-e29b-41d4-a716-446655440004",
      "email": "admin@berezhok.ru",
      "name": "Админ Админович",
      "role": "super_admin",
      "permissions": ["*"]
    }
  }
}
```

---

## 2. Клиентское API (мобильное приложение)

**Заголовок:** `Authorization: Bearer {token}`

### 2.1 Получить профиль

**Endpoint:** `GET /customer/profile`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "phone": "+79001234567",
    "name": "Иван",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

---

### 2.2 Обновить профиль

**Endpoint:** `PATCH /customer/profile`

**Request:**
```json
{
  "name": "Иван Петров"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "phone": "+79001234567",
    "name": "Иван Петров",
    "updated_at": "2025-01-20T14:22:00Z"
  }
}
```

---

### 2.3 Поиск заведений

**Endpoint:** `GET /customer/locations`

**Query параметры:**
- `lat` — широта (обязательно)
- `lng` — долгота (обязательно)
- `radius` — радиус поиска в метрах (по умолчанию 5000)
- `category` — код категории (опционально)
- `limit`, `offset` — пагинация

**Пример запроса:**
```
GET /customer/locations?lat=55.7558&lng=37.6173&radius=2000&category=bakery&limit=20
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "850e8400-e29b-41d4-a716-446655440003",
        "name": "Пекарня на Ленина",
        "category": {
          "code": "bakery",
          "name": "Пекарня",
          "icon_url": "/icons/bakery.svg"
        },
        "address": "ул. Ленина, 5",
        "distance": 850,
        "coordinates": {
          "lat": 55.7558,
          "lng": 37.6173
        },
        "rating": {
          "average": 4.5,
          "total_reviews": 120
        },
        "logo_url": "/uploads/logo_bakery.jpg",
        "active_boxes_count": 3
      }
    ],
    "pagination": {
      "total": 15,
      "limit": 20,
      "offset": 0,
      "has_more": false
    }
  }
}
```

---

### 2.4 Получить детали заведения

**Endpoint:** `GET /customer/locations/{location_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "850e8400-e29b-41d4-a716-446655440003",
    "name": "Пекарня на Ленина",
    "category": {
      "code": "bakery",
      "name": "Пекарня"
    },
    "address": "ул. Ленина, 5, Москва",
    "coordinates": {
      "lat": 55.7558,
      "lng": 37.6173
    },
    "phone": "+74951234567",
    "working_hours": {
      "mon": "08:00-22:00",
      "tue": "08:00-22:00",
      "wed": "08:00-22:00",
      "thu": "08:00-22:00",
      "fri": "08:00-23:00",
      "sat": "09:00-23:00",
      "sun": "09:00-21:00"
    },
    "logo_url": "/uploads/logo.jpg",
    "cover_image_url": "/uploads/cover.jpg",
    "gallery": ["/uploads/img1.jpg", "/uploads/img2.jpg"],
    "rating": {
      "average": 4.5,
      "total_reviews": 120,
      "distribution": {
        "5": 80,
        "4": 30,
        "3": 7,
        "2": 2,
        "1": 1
      }
    },
    "active_boxes": [
      {
        "id": "a50e8400-e29b-41d4-a716-446655440005",
        "name": "Вечерний сюрприз",
        "description": "Свежая выпечка дня: круассаны, булочки, пирожные",
        "original_price": 500,
        "discount_price": 199,
        "quantity_available": 5,
        "pickup_time": {
          "start": "18:00",
          "end": "19:00"
        },
        "image_url": "/uploads/box.jpg"
      }
    ]
  }
}
```

---

### 2.5 Создать заказ

**Endpoint:** `POST /customer/orders`

**Request:**
```json
{
  "box_id": "a50e8400-e29b-41d4-a716-446655440005"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "order_id": "b60e8400-e29b-41d4-a716-446655440006",
    "payment_url": "https://yookassa.ru/checkout/...",
    "amount": 199,
    "expires_at": "2025-01-20T15:00:00Z"
  }
}
```

**Errors:**
- `400` — бокс недоступен (sold out)
- `409` — у пользователя уже есть активный заказ на это время

---

### 2.6 Webhook от ЮKassa (внутренний)

**Endpoint:** `POST /webhooks/yookassa`

После успешной оплаты backend обновляет статус заказа и отправляет уведомление клиенту.

---

### 2.7 Получить детали заказа

**Endpoint:** `GET /customer/orders/{order_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "b60e8400-e29b-41d4-a716-446655440006",
    "status": "confirmed",
    "pickup_code": "AB12CD34",
    "qr_code_url": "/qr/AB12CD34.png",
    "amount": 199,
    "box": {
      "name": "Вечерний сюрприз",
      "image_url": "/uploads/box.jpg"
    },
    "location": {
      "name": "Пекарня на Ленина",
      "address": "ул. Ленина, 5",
      "phone": "+74951234567",
      "coordinates": {
        "lat": 55.7558,
        "lng": 37.6173
      }
    },
    "pickup_time": {
      "start": "2025-01-20T18:00:00Z",
      "end": "2025-01-20T19:00:00Z"
    },
    "created_at": "2025-01-20T14:30:00Z",
    "confirmed_at": "2025-01-20T14:32:00Z"
  }
}
```

---

### 2.8 Подтвердить получение заказа

**Endpoint:** `POST /customer/orders/{order_id}/confirm-pickup`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "order_id": "b60e8400-e29b-41d4-a716-446655440006",
    "status": "completed",
    "can_review": true
  }
}
```

---

### 2.9 Открыть спор

**Endpoint:** `POST /customer/orders/{order_id}/dispute`

**Request:**
```json
{
  "reason": "Не получил заказ, сотрудник не выдал"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "dispute_id": "c70e8400-e29b-41d4-a716-446655440007",
    "status": "open",
    "message": "Спор открыт. Мы рассмотрим его в течение 24 часов."
  }
}
```

**Errors:**
- `400` — время для открытия спора истекло (>15 минут после выдачи)

---

### 2.10 Оставить отзыв

**Endpoint:** `POST /customer/reviews`

**Request:**
```json
{
  "order_id": "b60e8400-e29b-41d4-a716-446655440006",
  "rating": 5,
  "comment": "Отличный бокс! Всё свежее и вкусное"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "review_id": "d80e8400-e29b-41d4-a716-446655440008",
    "rating": 5,
    "comment": "Отличный бокс! Всё свежее и вкусное",
    "created_at": "2025-01-20T20:15:00Z"
  }
}
```

**Errors:**
- `400` — заказ ещё не завершён
- `409` — отзыв уже оставлен

---

### 2.11 История заказов

**Endpoint:** `GET /customer/orders`

**Query параметры:**
- `status` — фильтр по статусу (опционально)
- `limit`, `offset` — пагинация

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "b60e8400-e29b-41d4-a716-446655440006",
        "status": "completed",
        "pickup_code": "AB12CD34",
        "amount": 199,
        "box_name": "Вечерний сюрприз",
        "location_name": "Пекарня на Ленина",
        "pickup_time_start": "2025-01-20T18:00:00Z",
        "created_at": "2025-01-20T14:30:00Z",
        "has_review": true
      }
    ],
    "pagination": {
      "total": 15,
      "limit": 20,
      "offset": 0,
      "has_more": false
    }
  }
}
```

---

### 2.12 Отзывы о заведении

**Endpoint:** `GET /customer/locations/{location_id}/reviews`

**Query параметры:**
- `limit`, `offset` — пагинация

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "d80e8400-e29b-41d4-a716-446655440008",
        "rating": 5,
        "comment": "Отличный бокс! Всё свежее и вкусное",
        "user_name": "Иван",
        "created_at": "2025-01-20T20:15:00Z"
      }
    ],
    "pagination": {
      "total": 120,
      "limit": 20,
      "offset": 0,
      "has_more": true
    }
  }
}
```

---

## 3. Partner API (веб-панель партнёров)

**Заголовок:** `Authorization: Bearer {token}`

### 3.1 Получить профиль партнёра

**Endpoint:** `GET /partner/profile`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "partner": {
      "id": "750e8400-e29b-41d4-a716-446655440002",
      "legal_name": "ООО Вкусная пекарня",
      "brand_name": "Мама печёт",
      "status": "active",
      "commission_rate": 0.10,
      "promo_until": "2025-04-15"
    },
    "employee": {
      "id": "650e8400-e29b-41d4-a716-446655440001",
      "name": "Иван Иванов",
      "email": "owner@bakery.ru",
      "role": "owner"
    },
    "location": {
      "id": "850e8400-e29b-41d4-a716-446655440003",
      "name": "Пекарня на Ленина",
      "address": "ул. Ленина, 5"
    }
  }
}
```

---

### 3.2 Список сюрприз-боксов

**Endpoint:** `GET /partner/boxes`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "a50e8400-e29b-41d4-a716-446655440005",
        "name": "Вечерний сюрприз",
        "description": "Свежая выпечка дня",
        "original_price": 500,
        "discount_price": 199,
        "quantity_available": 5,
        "pickup_time": {
          "start": "18:00",
          "end": "19:00"
        },
        "status": "active",
        "created_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

---

### 3.3 Создать сюрприз-бокс

**Endpoint:** `POST /partner/boxes`

**Request:**
```json
{
  "name": "Утренний бокс",
  "description": "Круассаны, булочки, кофе",
  "original_price": 400,
  "discount_price": 149,
  "quantity_available": 10,
  "pickup_time_start": "08:00",
  "pickup_time_end": "09:00",
  "image_url": "/uploads/morning_box.jpg"
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "a60e8400-e29b-41d4-a716-446655440009",
    "name": "Утренний бокс",
    "status": "inactive",
    "message": "Бокс создан. Активируйте его для публикации."
  }
}
```

**Errors:**
- `400` — у партнёра уже 5 активных боксов (лимит)

---

### 3.4 Обновить сюрприз-бокс

**Endpoint:** `PATCH /partner/boxes/{box_id}`

**Request:**
```json
{
  "quantity_available": 15,
  "status": "active"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "a60e8400-e29b-41d4-a716-446655440009",
    "quantity_available": 15,
    "status": "active",
    "updated_at": "2025-01-20T16:00:00Z"
  }
}
```

---

### 3.5 Заказы требующие подтверждения

**Endpoint:** `GET /partner/orders/pending-confirmation`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "b60e8400-e29b-41d4-a716-446655440006",
        "pickup_code": "AB12CD34",
        "box_name": "Вечерний сюрприз",
        "amount": 199,
        "customer_phone": "+7900***4567",
        "pickup_time_start": "2025-01-20T18:00:00Z",
        "confirmation_deadline": "2025-01-20T15:02:00Z",
        "time_left_seconds": 450,
        "created_at": "2025-01-20T14:32:00Z"
      }
    ]
  }
}
```

---

### 3.6 Подтвердить заказ

**Endpoint:** `POST /partner/orders/{order_id}/confirm`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "order_id": "b60e8400-e29b-41d4-a716-446655440006",
    "status": "confirmed",
    "message": "Заказ подтверждён. Клиент получил уведомление."
  }
}
```

**Errors:**
- `400` — время подтверждения истекло
- `409` — заказ уже подтверждён/отменён

---

### 3.7 Поиск заказа по коду

**Endpoint:** `GET /partner/orders/by-code/{pickup_code}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "b60e8400-e29b-41d4-a716-446655440006",
    "pickup_code": "AB12CD34",
    "status": "confirmed",
    "box_name": "Вечерний сюрприз",
    "customer_phone": "+7900***4567",
    "pickup_time": {
      "start": "2025-01-20T18:00:00Z",
      "end": "2025-01-20T19:00:00Z"
    }
  }
}
```

**Errors:**
- `404` — заказ не найден

---

### 3.8 Выдать заказ

**Endpoint:** `POST /partner/orders/{order_id}/pickup`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "order_id": "b60e8400-e29b-41d4-a716-446655440006",
    "status": "picked_up",
    "message": "Заказ отмечен как выданный. Ожидаем подтверждения от клиента."
  }
}
```

---

### 3.9 Статистика

**Endpoint:** `GET /partner/stats`

**Query параметры:**
- `period` — день/неделя/месяц (по умолчанию неделя)

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "period": "week",
    "orders": {
      "total": 45,
      "completed": 42,
      "cancelled": 3
    },
    "revenue": {
      "gross": 8550,
      "commission": 855,
      "net": 7695
    },
    "rating": {
      "average": 4.5,
      "total_reviews": 38
    },
    "missed_confirmations": {
      "count": 1,
      "percentage": 2.2
    }
  }
}
```

---

### 3.10 История выплат

**Endpoint:** `GET /partner/payouts`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "e90e8400-e29b-41d4-a716-446655440010",
        "period": {
          "start": "2025-01-13",
          "end": "2025-01-19"
        },
        "gross_amount": 8550,
        "commission_rate": 0.10,
        "commission_amount": 855,
        "net_amount": 7695,
        "orders_count": 45,
        "status": "completed",
        "processed_at": "2025-01-20T10:00:00Z"
      }
    ]
  }
}
```

---

## 4. Admin API (админ-панель)

**Заголовок:** `Authorization: Bearer {token}`  
**Требуется:** роль `admin` или выше

### 4.1 Заявки на регистрацию

**Endpoint:** `GET /admin/applications`

**Query параметры:**
- `status` — pending/approved/rejected

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "f00e8400-e29b-41d4-a716-446655440011",
        "business_name": "Пекарня Ромашка",
        "contact_name": "Пётр Петров",
        "contact_email": "petr@romashka.ru",
        "contact_phone": "+79001234567",
        "category": "bakery",
        "address": "ул. Пушкина, 10",
        "status": "pending",
        "created_at": "2025-01-19T12:00:00Z"
      }
    ]
  }
}
```

---

### 4.2 Одобрить заявку

**Endpoint:** `POST /admin/applications/{application_id}/approve`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "application_id": "f00e8400-e29b-41d4-a716-446655440011",
    "partner_id": "750e8400-e29b-41d4-a716-446655440012",
    "owner_email": "petr@romashka.ru",
    "message": "Партнёр создан. На email отправлены учётные данные."
  }
}
```

---

### 4.3 Отклонить заявку

**Endpoint:** `POST /admin/applications/{application_id}/reject`

**Request:**
```json
{
  "reason": "Недостаточно информации о бизнесе"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "application_id": "f00e8400-e29b-41d4-a716-446655440011",
    "status": "rejected"
  }
}
```

---

### 4.4 Список партнёров

**Endpoint:** `GET /admin/partners`

**Query параметры:**
- `status` — фильтр по статусу
- `search` — поиск по названию
- `limit`, `offset` — пагинация

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "750e8400-e29b-41d4-a716-446655440002",
        "legal_name": "ООО Вкусная пекарня",
        "brand_name": "Мама печёт",
        "status": "active",
        "locations_count": 1,
        "total_orders": 156,
        "total_revenue": 29640,
        "missed_confirmation_rate": 2.5,
        "created_at": "2024-12-01T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 45,
      "limit": 20,
      "offset": 0,
      "has_more": true
    }
  }
}
```

---

### 4.5 Детали партнёра

**Endpoint:** `GET /admin/partners/{partner_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440002",
    "legal_name": "ООО Вкусная пекарня",
    "brand_name": "Мама печёт",
    "status": "active",
    "commission_rate": 0.10,
    "promo_until": "2025-04-15",
    "legal_info": {
      "inn": "1234567890",
      "ogrn": "1234567890123",
      "legal_address": "г. Москва, ул. Ленина, 1"
    },
    "locations": [
      {
        "id": "850e8400-e29b-41d4-a716-446655440003",
        "name": "Пекарня на Ленина",
        "address": "ул. Ленина, 5",
        "status": "active"
      }
    ],
    "stats": {
      "total_orders": 156,
      "total_revenue": 29640,
      "missed_confirmations": 4,
      "missed_confirmation_rate": 2.56
    },
    "created_at": "2024-12-01T10:00:00Z"
  }
}
```

---

### 4.6 Заблокировать партнёра

**Endpoint:** `POST /admin/partners/{partner_id}/block`

**Request:**
```json
{
  "reason": "Множественные жалобы клиентов"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "partner_id": "750e8400-e29b-41d4-a716-446655440002",
    "status": "blocked"
  }
}
```

---

### 4.7 Список споров

**Endpoint:** `GET /admin/disputes`

**Query параметры:**
- `status` — open/under_review/resolved

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "c70e8400-e29b-41d4-a716-446655440007",
        "order_id": "b60e8400-e29b-41d4-a716-446655440006",
        "initiated_by": "user",
        "reason": "Не получил заказ",
        "status": "open",
        "customer": {
          "phone": "+7900***4567",
          "name": "Иван"
        },
        "partner": {
          "name": "Пекарня на Ленина"
        },
        "created_at": "2025-01-20T18:20:00Z"
      }
    ]
  }
}
```

---

### 4.8 Детали спора

**Endpoint:** `GET /admin/disputes/{dispute_id}`

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "c70e8400-e29b-41d4-a716-446655440007",
    "order": {
      "id": "b60e8400-e29b-41d4-a716-446655440006",
      "pickup_code": "AB12CD34",
      "amount": 199,
      "status": "disputed"
    },
    "customer": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "phone": "+79001234567",
      "name": "Иван"
    },
    "partner": {
      "id": "750e8400-e29b-41d4-a716-446655440002",
      "name": "ООО Вкусная пекарня"
    },
    "location": {
      "name": "Пекарня на Ленина"
    },
    "reason": "Не получил заказ",
    "status": "open",
    "timeline": [
      {
        "timestamp": "2025-01-20T18:15:00Z",
        "action": "employee_scanned_code",
        "actor": "employee@bakery.ru"
      },
      {
        "timestamp": "2025-01-20T18:20:00Z",
        "action": "user_opened_dispute",
        "actor": "+79001234567"
      }
    ],
    "created_at": "2025-01-20T18:20:00Z"
  }
}
```

---

### 4.9 Решить спор

**Endpoint:** `POST /admin/disputes/{dispute_id}/resolve`

**Request:**
```json
{
  "resolution": "resolved_for_user",
  "comment": "Сотрудник подтвердил ошибку, возвращаем деньги клиенту"
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "dispute_id": "c70e8400-e29b-41d4-a716-446655440007",
    "resolution": "resolved_for_user",
    "order_status": "refunded",
    "message": "Спор решён в пользу клиента. Возврат средств инициирован."
  }
}
```

---

### 4.10 Статистика платформы

**Endpoint:** `GET /admin/stats`

**Query параметры:**
- `period` — day/week/month/year

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "period": "month",
    "orders": {
      "total": 4520,
      "completed": 4320,
      "cancelled": 150,
      "disputed": 50
    },
    "revenue": {
      "gmv": 856800,
      "commission": 85680,
      "net_revenue": 85680
    },
    "users": {
      "total": 12450,
      "new": 850,
      "active": 3200
    },
    "partners": {
      "total": 45,
      "active": 42,
      "new": 3
    },
    "top_partners": [
      {
        "id": "750e8400-e29b-41d4-a716-446655440002",
        "name": "ООО Вкусная пекарня",
        "orders": 320,
        "revenue": 60800
      }
    ]
  }
}
```

---

## 5. Webhooks

### 5.1 ЮKassa webhook (payment success)

**Endpoint:** `POST /webhooks/yookassa`

**Request (от ЮKassa):**
```json
{
  "type": "notification",
  "event": "payment.succeeded",
  "object": {
    "id": "external_payment_id_123",
    "status": "succeeded",
    "amount": {
      "value": "199.00",
      "currency": "RUB"
    },
    "metadata": {
      "order_id": "b60e8400-e29b-41d4-a716-446655440006"
    }
  }
}
```

**Backend обрабатывает:**
1. Обновляет статус заказа: `pending` → `paid`
2. Устанавливает дедлайн подтверждения (30 минут)
3. Отправляет push-уведомление клиенту
4. Отправляет уведомление партнёру

---

## 6. Коды ошибок

| Код | Описание |
|-----|----------|
| `VALIDATION_ERROR` | Ошибка валидации данных |
| `UNAUTHORIZED` | Не авторизован |
| `FORBIDDEN` | Нет прав доступа |
| `NOT_FOUND` | Ресурс не найден |
| `CONFLICT` | Конфликт (дубликат, уже существует) |
| `PAYMENT_FAILED` | Ошибка оплаты |
| `BOX_UNAVAILABLE` | Бокс недоступен |
| `CONFIRMATION_EXPIRED` | Время подтверждения истекло |
| `DISPUTE_TIME_EXPIRED` | Время для открытия спора истекло |
| `RATE_LIMIT_EXCEEDED` | Превышен лимит запросов |

---

## 7. Rate Limiting

- **Общий лимит:** 100 запросов/минуту на IP
- **Аутентификация:** 5 попыток/5 минут на IP
- **SMS коды:** 3 запроса/час на номер телефона

**Заголовки ответа при rate limit:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642680000
```

**При превышении лимита:** `429 Too Many Requests`
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Слишком много запросов. Повторите через 60 секунд.",
    "retry_after": 60
  }
}
```

---

**Конец документа**