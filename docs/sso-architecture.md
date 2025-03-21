# Архитектура SSO

```mermaid
graph TB
    subgraph External["Внешний слой (gRPC API)"]
        Client[gRPC Клиент]
        Gateway[gRPC Сервер]
        
        Client -->|gRPC| Gateway
        Gateway -->|Login| AuthService
        Gateway -->|Register| AuthService
        Gateway -->|IsAdmin| AuthService
    end
    
    subgraph Business["Слой бизнес-логики (Services)"]
        AuthService[Auth Сервис]
        JWT[JWT Модуль]
        
        AuthService -->|Генерация| JWT
        AuthService -->|Хеширование| AuthService
    end
    
    subgraph Data["Слой данных (Storage)"]
        SQLite[(SQLite DB)]
        Tables[Таблицы]
        
        subgraph Tables
            Users[users]
            Apps[apps]
        end
        
        SQLite -->|Хранит| Users
        SQLite -->|Хранит| Apps
    end
    
    subgraph Utils["Вспомогательные компоненты"]
        Logger[Logger]
        Config[Config]
    end
    
    AuthService -->|CRUD| SQLite
    AuthService -->|Логирование| Logger
    AuthService -->|Конфигурация| Config
    
    classDef external fill:#f9f,stroke:#333,stroke-width:2px;
    classDef business fill:#bbf,stroke:#333,stroke-width:2px;
    classDef data fill:#bfb,stroke:#333,stroke-width:2px;
    classDef utils fill:#fbb,stroke:#333,stroke-width:2px;
    
    class External external;
    class Business business;
    class Data data;
    class Utils utils;
```

## Описание компонентов

### Внешний слой (gRPC API)
- **gRPC Клиент**: Внешние приложения, использующие SSO
- **gRPC Сервер**: Обработка входящих запросов и валидация данных

### Слой бизнес-логики (Services)
- **Auth Сервис**: Основная бизнес-логика аутентификации и авторизации
- **JWT Модуль**: Генерация и валидация токенов

### Слой данных (Storage)
- **SQLite DB**: Хранилище данных
- **Таблицы**:
  - `users`: Пользователи системы
  - `apps`: Зарегистрированные приложения

### Вспомогательные компоненты
- **Logger**: Логирование операций
- **Config**: Управление конфигурацией

## Процессы

### Аутентификация
1. Клиент отправляет запрос на Login
2. gRPC сервер валидирует данные
3. Auth сервис проверяет пользователя
4. Генерируется JWT токен
5. Токен возвращается клиенту

### Регистрация
1. Клиент отправляет запрос на Register
2. gRPC сервер валидирует данные
3. Auth сервис хеширует пароль
4. Данные сохраняются в SQLite
5. ID пользователя возвращается клиенту

### Проверка прав
1. Клиент отправляет запрос на IsAdmin
2. gRPC сервер валидирует данные
3. Проверяются права в SQLite
4. Результат возвращается клиенту 