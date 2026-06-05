# Бережок

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16_+_PostGIS-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/RabbitMQ-FF6600?logo=rabbitmq&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

Платформа для фудшеринга по модели сюрприз-боксов: заведения (пекарни, кафе, рестораны, магазины) продают нераспроданную за день еду со скидкой, клиенты её выкупают и забирают в заданное окно. По смыслу - российский аналог Too Good To Go.

Репозиторий — это монорепа: Go-бэкенд (модульный монолит) + React-фронтенд для партнёрской и админской панелей.

## Стек

**Бэкенд** — Go, chi, pgx, PostgreSQL 16 + PostGIS, Redis, RabbitMQ, S3 (Yandex Object Storage), JWT. SQL генерируется через sqlc, миграции - golang-migrate.

**Фронтенд** — React 19 + Vite. MobX для состояния, Tailwind, axios.

**Инфра** — Docker Compose, Traefik (TLS через Let's Encrypt), Taskfile + Makefile для рутины.

Платежи — ЮKassa. Геокодинг — Yandex Geocoder.

## Структура

```
cmd/
  api/              # основной HTTP-сервис
  payout-worker/    # воркер выплат (расчёт / отправка / опрос статусов)
  admin-bootstrap/  # создание первого администратора
  seed/             # наполнение БД тестовыми данными
internal/
  adapters/         # внешние интеграции: postgres, redis, s3, sms, yookassa
  modules/          # фичевые модули (см. ниже)
  shared/           # auth, jwt, config, middleware, общие домены и ошибки
  lib/              # логгер, валидатор, конвертеры pgx-типов
migrations/         # SQL-миграции
frontend/           # React-панели партнёра и админа
docs/               # подробная спецификация продукта и контракт API
```

Каждый модуль в `internal/modules/` устроен одинаково: `domain` (бизнес-сущности), `repository` (доступ к данным поверх sqlc), `service` (логика), `handlers` + `handlers/dto` (HTTP и DTO), `errors` (sentinel-ошибки модуля). Модули: `auth`, `customer`, `partner`, `admin`, `catalog`, `order`, `payment`, `payout`, `review`, `eco`, `media`.

## Как это работает

Центральная сущность — заказ, и почти вся логика крутится вокруг его статусов. Клиент бронирует и оплачивает бокс через ЮKassa, дальше заказ проходит цепочку: `paid` → партнёр подтверждает (`confirmed`) → сотрудник выдаёт по коду (`picked_up`) → автозавершение (`completed`). На неподтверждённых вовремя заказах и спорах ветки расходятся — авто-отмена с возвратом или `disputed`. Выручка партнёрам выплачивается не сразу, а пачками: `payout-worker` раз в период считает completed-заказы, удерживает комиссию и отправляет выплату.

Фронтенд ходит в API через единый axios-клиент; все ответы завёрнуты в конверт `{ "success": bool, "data": ... }` (или `error` с кодом и сообщением). Аутентификация трёх видов: партнёр (email + пароль), клиент (телефон + SMS-код), админ — роль проверяется в middleware.

Подробное описание бизнес-процессов (жизненный цикл заказа, выплаты, споры, роли) лежит в [docs/berezhok.md](docs/berezhok.md). Контракт API — в [docs/api_contract.md](docs/api_contract.md). Для агентов и общих соглашений по коду — [AGENTS.md](AGENTS.md).

## Команды

Бэкенд (Taskfile):

```bash
task lint            # golangci-lint
task format          # gofumpt + gci (сортировка импортов)
task tests           # go test ./internal/...
task e2e_tests       # E2E-тесты API (нужна БД)
```

Makefile:

```bash
make sql-gen         # sqlc generate
make migrate-create name=foo
make migrate-up / make migrate-down
make pre-commit      # все линтеры и форматтеры
```

Воркер выплат запускается по режимам:

```bash
make pay-calc        # сформировать выплаты за период
make pay-disp        # отправить выплаты в ЮKassa
make pay-poll        # опросить статусы
```

Фронтенд:

```bash
npm run dev / npm run build / npm run lint
```
