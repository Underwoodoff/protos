# Архитектура слоев SSO

## Схема слоев

```mermaid
graph TB
    subgraph "Внешний слой"
        Client[Клиент]
        gRPC[gRPC API]
        HTTP[HTTP API]
        Middleware[Middleware]
        Handlers[Handlers]
    end

    subgraph "Бизнес-логика"
        Auth[Auth Service]
        User[User Service]
        App[App Service]
        Domain[Domain Models]
    end

    subgraph "Слой данных"
        Storage[Storage Interface]
        SQLite[(SQLite)]
    end

    subgraph "Вспомогательные компоненты"
        Config[Config]
        Logger[Logger]
        JWT[JWT]
        Metrics[Metrics]
    end

    %% Внешние связи
    Client -->|gRPC| gRPC
    Client -->|HTTP/2| HTTP
    HTTP -->|Middleware Chain| Middleware
    Middleware -->|Request/Response| Handlers
    gRPC -->|Service Calls| Auth
    Handlers -->|Service Calls| Auth

    %% Внутренние связи
    Auth -->|User Management| User
    Auth -->|App Validation| App
    Auth -->|Domain Logic| Domain
    User -->|Data Access| Storage
    App -->|Data Access| Storage
    Storage -->|CRUD| SQLite

    %% Вспомогательные связи
    Config -->|Configuration| Auth
    Config -->|Configuration| Storage
    Logger -->|Logging| Auth
    Logger -->|Logging| Handlers
    JWT -->|Token Generation| Auth
    Metrics -->|Monitoring| Auth
    Metrics -->|Monitoring| Handlers

    %% Стили
    classDef external fill:#f9f,stroke:#333,stroke-width:2px;
    classDef business fill:#bbf,stroke:#333,stroke-width:2px;
    classDef data fill:#bfb,stroke:#333,stroke-width:2px;
    classDef utils fill:#fbb,stroke:#333,stroke-width:2px;

    class Client,gRPC,HTTP,Middleware,Handlers external;
    class Auth,User,App,Domain business;
    class Storage,SQLite data;
    class Config,Logger,JWT,Metrics utils;
```

## Описание слоев

### 1. Внешний слой (External Layer)
- **gRPC API**
  - Обработка gRPC запросов
  - Валидация входных данных
  - Сериализация/десериализация сообщений

- **HTTP API**
  - Обработка HTTP/2 запросов
  - RESTful API для аутентификации
  - Маршрутизация запросов

- **Middleware**
  - Аутентификация
  - Логирование
  - Метрики
  - Обработка ошибок
  - Rate limiting

- **Handlers**
  - Обработчики HTTP запросов
  - Преобразование HTTP в доменные модели
  - Валидация бизнес-правил

### 2. Бизнес-логика (Business Logic)
- **Auth Service**
  - Аутентификация пользователей
  - Генерация JWT токенов
  - Валидация токенов
  - Управление сессиями

- **User Service**
  - Управление пользователями
  - Регистрация
  - Проверка прав доступа
  - Управление профилями

- **App Service**
  - Управление приложениями
  - Валидация приложений
  - Управление разрешениями
  - Конфигурация приложений

- **Domain Models**
  - Бизнес-сущности (User, App, Token)
  - Валидация бизнес-правил
  - Доменные события

### 3. Слой данных (Data Layer)
- **Storage Interface**
  - Абстракция доступа к данным
  - Repository pattern
  - Unit of Work
  - Транзакции

- **SQLite**
  - Постоянное хранение данных
  - Транзакционность
  - Миграции схемы
  - Оптимизация запросов

### 4. Вспомогательные компоненты (Utilities)
- **Config**
  - Конфигурация приложения
  - Переменные окружения
  - Feature flags
  - Настройки компонентов

- **Logger**
  - Структурированное логирование
  - Уровни логирования
  - Контекст
  - Трассировка запросов

- **JWT**
  - Генерация токенов
  - Валидация токенов
  - Управление ключами
  - Обновление токенов

- **Metrics**
  - Prometheus метрики
  - Мониторинг
  - Трейсинг
  - Профилирование

## Принципы взаимодействия

1. **Зависимости направлены внутрь**
   - Внешний слой зависит от бизнес-логики
   - Бизнес-логика зависит от абстракций слоя данных
   - Вспомогательные компоненты доступны всем слоям

2. **Инверсия зависимостей**
   - Слои зависят от абстракций
   - Интерфейсы определяются потребителем
   - Внедрение зависимостей через конструкторы

3. **Чистая архитектура**
   - Разделение ответственности
   - Независимость от фреймворков
   - Тестируемость компонентов

## Преимущества архитектуры

1. **Масштабируемость**
   - Легко добавлять новые функции
   - Возможность горизонтального масштабирования
   - Независимость компонентов

2. **Поддерживаемость**
   - Четкое разделение ответственности
   - Легкость внесения изменений
   - Понятная структура кода

3. **Тестируемость**
   - Изолированные компоненты
   - Легкость написания unit-тестов
   - Возможность мокирования зависимостей

4. **Безопасность**
   - Контроль доступа на уровне middleware
   - Защита от перегрузок
   - Валидация на всех уровнях 