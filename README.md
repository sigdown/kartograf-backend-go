## Обзор

Kartograf — клиент-серверная система для наложения пользовательских и исторических карт поверх современной картографической подложки.

Основные сценарии:

- просмотр публичного каталога карт
- открытие карточки карты
- показ тайлов карты на клиенте
- хранение и управление удалёнными пользовательскими точками
- загрузка и замена PMTiles-архивов через админский сценарий

## Что есть в этом репозитории

Этот репозиторий содержит backend API.

Текущая зона ответственности backend:

- JWT-аутентификация
- публичный каталог карт
- получение карты по `slug`
- presigned download для архива
- CRUD удалённых точек для авторизованного пользователя
- админское создание карты и замена архива через presigned upload

## Правила хранения карт

- `slug` карты уникален
- `slug` карты не меняется после создания
- активный PMTiles-объект хранится по ключу `kartograf/<slug>.pmtiles`
- загрузка архива идёт напрямую в object storage через presigned URL
- скачивание архива возвращается как presigned URL

## Настройка окружения

Создай локальный `.env` на основе примера:

```bash
cp .env.example .env
```

Заполни основные переменные:

- `APP_HOST`, `APP_PORT` — адрес и порт API
- `POSTGRES_DSN` — строка подключения к PostgreSQL
- `S3_ENDPOINT` — endpoint object storage
- `S3_REGION` — регион object storage
- `S3_ACCESS_KEY`, `S3_SECRET_KEY` — credentials для object storage
- `S3_BUCKET` — bucket для PMTiles-архивов
- `S3_USE_PATH_STYLE` — нужен ли path-style доступ
- `AUTH_JWT_SECRET` — секрет для подписи JWT
- `AUTH_ACCESS_TOKEN_TTL` — время жизни access token

Для локальной разработки можно использовать значения из `.env.example`.

## Локальная разработка

Большая часть рутинных команд вынесена в `Makefile`.

Для локальной инфраструктуры используется `docker compose`.
Локальный compose поднимает:

- `postgres`
- `minio`

Можно работать и прямыми командами:

```bash
docker compose up -d
docker compose down
```

Но в обычной работе удобнее использовать цели из `Makefile`.

Поднять инфраструктуру:

```bash
make db-up
```

Посмотреть логи PostgreSQL:

```bash
make db-logs
```

Применить миграции:

```bash
make migrate-up
```

Запустить API:

```bash
make run
```

Прогнать тесты:

```bash
make test
```

Проверить сборку:

```bash
make build
```

Остановить инфраструктуру:

```bash
make db-down
```

## Документация

- контракт API: [docs/api.md](docs/api.md)
