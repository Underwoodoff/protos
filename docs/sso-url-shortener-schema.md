# Схема взаимодействия SSO и URL Shortener

```mermaid
graph LR
    Client[Клиент] -->|HTTP REST| URLShortener[URL Shortener]
    URLShortener -->|gRPC| SSO[SSO Service]
    
    subgraph URL Shortener
        REST[HTTP Server]
        AuthMiddleware[Auth Middleware]
        URLService[URL Service]
    end
    
    subgraph SSO Service
        AuthService[Auth Service]
        Storage[(SQLite DB)]
    end
    
    Client -->|1. HTTP POST /api/urls| REST
    REST -->|2. Проверка токена| AuthMiddleware
    AuthMiddleware -->|3. gRPC IsAdmin()| AuthService
    AuthService -->|4. Проверка в БД| Storage
    
    AuthMiddleware -->|5. Разрешение доступа| URLService
    URLService -->|6. HTTP Response| Client
    
    style Client fill:#f9f,stroke:#333,stroke-width:2px
    style URLShortener fill:#bbf,stroke:#333,stroke-width:2px
    style SSO fill:#bfb,stroke:#333,stroke-width:2px
    style AuthMiddleware fill:#fbb,stroke:#333,stroke-width:2px
    style AuthService fill:#fbb,stroke:#333,stroke-width:2px
    style Storage fill:#fbb,stroke:#333,stroke-width:2px
```

## Описание взаимодействия

1. **Клиент -> URL Shortener (HTTP REST)**
   - Протокол: HTTP
   - Метод: POST
   - Путь: /api/urls
   - Заголовок: Authorization: Bearer <jwt_token>

2. **URL Shortener -> Auth Middleware**
   - Проверка наличия и валидности JWT токена
   - Извлечение userID из токена

3. **Auth Middleware -> SSO Service (gRPC)**
   - Протокол: gRPC
   - Метод: IsAdmin()
   - Параметры: userID
   - Интерцепторы: grpclog, grpcretry

4. **SSO Service -> Storage**
   - Проверка прав администратора в SQLite
   - Возврат результата

5. **Auth Middleware -> URL Service**
   - Разрешение/запрет доступа на основе результата проверки
   - Добавление информации о пользователе в контекст

6. **URL Service -> Клиент (HTTP REST)**
   - Успешный ответ или ошибка авторизации
   - Результат операции с URL 