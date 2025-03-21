# Архитектура системы

## Общая схема взаимодействия

```mermaid
graph TB
    Client[Клиент]
    SSO[SSO Service]
    URLShortener[URL Shortener Service]
    DB[(PostgreSQL)]
    Redis[(Redis)]

    subgraph "SSO Service"
        SSO --> Auth[Auth Service]
        SSO --> UserStorage[(User Storage)]
    end

    subgraph "URL Shortener Service"
        URLShortener --> URLService[URL Service]
        URLShortener --> URLStorage[(URL Storage)]
    end

    %% Внешние HTTP запросы
    Client -->|HTTP POST /api/v1/auth/login| SSO
    Client -->|HTTP POST /api/v1/auth/register| SSO
    Client -->|HTTP POST /api/v1/urls| URLShortener
    Client -->|HTTP GET /{short_url}| URLShortener

    %% Внутренние gRPC вызовы
    SSO -->|gRPC: ValidateToken| Auth
    Auth -->|gRPC: GetUser| UserStorage
    
    %% Внутренние вызовы URL Shortener
    URLShortener -->|gRPC: CreateURL| URLService
    URLService -->|gRPC: SaveURL| URLStorage
    URLService -->|Redis SET| Redis
    
    %% Внутренние вызовы при редиректе
    URLShortener -->|Redis GET| Redis
    URLShortener -->|gRPC: GetURL| URLStorage

    %% Стили для линий
    classDef http fill:#f9f,stroke:#333,stroke-width:2px;
    classDef grpc fill:#bbf,stroke:#333,stroke-width:2px;
    classDef redis fill:#bfb,stroke:#333,stroke-width:2px;
    
    class 1,2,3,4 http;
    class 5,6,7,8,9 grpc;
    class 10,11 redis;
```

## Описание взаимодействия

1. **Аутентификация пользователя**:
   - Клиент отправляет HTTP POST запрос на `/api/v1/auth/login` в SSO Service
   - SSO Service проверяет учетные данные через gRPC вызов Auth Service
   - Auth Service валидирует данные через gRPC вызов UserStorage
   - Возвращается JWT токен клиенту

2. **Создание короткой ссылки**:
   - Клиент отправляет HTTP POST запрос на `/api/v1/urls` в URL Shortener Service
   - URL Shortener Service через gRPC вызывает URL Service
   - URL Service сохраняет URL через gRPC в URLStorage
   - URL кэшируется в Redis через SET операцию

3. **Переход по короткой ссылке**:
   - Клиент отправляет HTTP GET запрос на `/{short_url}`
   - URL Shortener Service проверяет кэш через Redis GET
   - Если URL не найден в кэше, получает его через gRPC из URLStorage
   - Выполняется HTTP редирект на оригинальный URL

## Компоненты системы

### SSO Service
- Аутентификация пользователей (HTTP API)
- Управление сессиями (JWT)
- Генерация JWT токенов (golang-jwt)
- Хранение пользовательских данных (PostgreSQL)

### URL Shortener Service
- Генерация коротких ссылок (nanoid)
- Хранение соответствия коротких и оригинальных URL (PostgreSQL)
- Кэширование популярных ссылок (Redis)
- Перенаправление пользователей (HTTP)

### Хранилище данных
- PostgreSQL для хранения пользователей и URL
- Redis для кэширования коротких ссылок

### Протоколы и технологии
- HTTP/2 для внешнего API
- gRPC для внутренней коммуникации сервисов
- JWT для аутентификации
- PostgreSQL для постоянного хранения
- Redis для кэширования 