# Интеграция SSO с URL Shortener

В современной архитектуре микросервисов gRPC часто используется для взаимодействия между сервисами. В нашем случае мы рассмотрим интеграцию SSO с сервисом URL Shortener, который будет использовать аутентификацию и авторизацию через SSO. Это позволит нам продемонстрировать реальный пример взаимодействия между сервисами.

## Новая функциональность: метод IsAdmin()

Для интеграции с URL Shortener нам потребуется новый метод `IsAdmin()`, который будет определять, является ли пользователь администратором. Существующие методы `Login` и `RegisterNewUser` не подходят для этой задачи, так как они предназначены для прямого взаимодействия с клиентскими приложениями. В нашей схеме пользователь сначала получает токен от SSO, а затем использует его в URL Shortener.

Реализация метода `IsAdmin()` начинается с добавления нового поля в таблицу пользователей:

```sql
ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;
```

Затем мы добавляем метод в интерфейс Auth:

```go
type Auth interface {
    // ... существующие методы ...
    IsAdmin(ctx context.Context, userID int64) (bool, error)
}
```

И реализуем его в сервисе:

```go
func (s *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
    op := "AuthService.IsAdmin"
    s.log.Info("checking if user is admin", slog.Int64("user_id", userID))

    isAdmin, err := s.storage.IsAdmin(ctx, userID)
    if err != nil {
        s.log.Error("failed to check if user is admin", 
            slog.Int64("user_id", userID),
            slog.String("op", op),
            slog.Attr{
                Key:   "err",
                Value: slog.StringValue(err.Error()),
            })
        return false, fmt.Errorf("%s: %w", op, err)
    }

    s.log.Info("user is admin", 
        slog.Int64("user_id", userID),
        slog.Bool("is_admin", isAdmin))
    return isAdmin, nil
}
```

## Упрощения в реализации

В рамках текущей реализации мы сознательно идём на некоторые упрощения. Во-первых, метод `IsAdmin()` не защищён дополнительными механизмами безопасности. В идеале сервисы должны взаимодействовать через приватную сеть, что предоставляется многими облачными провайдерами, такими как Selectel.

Для более безопасной реализации можно использовать JWT-токены для авторизации сервисных запросов. Вместо передачи `userID` можно отправлять весь JWT-токен пользователя, что усложняет задачу для потенциальных злоумышленников. Хотя это не самый надёжный вариант, он прост в реализации и может быть использован как временное решение.

## Хранение ролей

В текущей реализации мы используем упрощённый подход к хранению ролей, добавляя колонку `is_admin` в таблицу `users`. Это решение не является оптимальным с точки зрения нормализации данных. В реальном проекте рекомендуется вынести роли и права в отдельную таблицу, особенно если у вас большое количество пользователей, но только небольшой процент из них являются администраторами.

Вот пример более правильной структуры для хранения ролей:

```sql
CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);
```

## Интеграция с URL Shortener

URL Shortener использует метод `IsAdmin()` для определения прав пользователя. Например, при попытке редактирования чужой записи:

```go
func (s *URLShortener) EditURL(ctx context.Context, userID int64, urlID string, newURL string) error {
    isAdmin, err := s.authClient.IsAdmin(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to check admin status: %w", err)
    }

    if !isAdmin {
        return errors.New("user is not authorized to edit this URL")
    }

    // ... логика редактирования URL ...
}
```

### Создание gRPC-клиента

Для взаимодействия с SSO сервисом в URL Shortener создаётся gRPC-клиент:

```go
grpcClient := ssov1.NewAuthClient(cc)
```

К клиенту подключаются два интерсептора:
- `grpclog` для логирования запросов и ответов
- `grpcretry` для автоматических повторных попыток при сбоях

### Middleware для авторизации

Поскольку URL Shortener является REST API сервисом, авторизация реализована в виде middleware. Это позволяет централизованно проверять права доступа для всех запросов:

```go
func extractBearerToken(r *http.Request) string {
    authHeader := r.Header.Get("Authorization")
    splitToken := strings.Split(authHeader, "Bearer ")
    if len(splitToken) != 2 {
        return ""
    }
    return splitToken[1]
}

type AuthMiddleware struct {
    ssoClient ssov1.AuthClient
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractBearerToken(r)
        if token == "" {
            http.Error(w, "Authorization token is required", http.StatusUnauthorized)
            return
        }

        claims, err := parseToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        isAdmin, err := m.ssoClient.IsAdmin(r.Context(), &ssov1.IsAdminRequest{
            UserId: claims.UserID,
        })
        if err != nil {
            http.Error(w, "Failed to check admin status", http.StatusInternalServerError)
            return
        }

        ctx := context.WithValue(r.Context(), "user", claims)
        if isAdmin {
            ctx = context.WithValue(ctx, "is_admin", true)
        }

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Заключение

Интеграция SSO с URL Shortener демонстрирует эффективное использование gRPC для взаимодействия между микросервисами. Реализованный подход обеспечивает:

1. Централизованную аутентификацию и авторизацию через SSO
2. Гибкое управление правами доступа через middleware
3. Безопасную передачу токенов между сервисами
4. Возможность масштабирования системы

В будущих версиях мы планируем расширить эту интеграцию, добавив:
- API Gateway для управления трафиком
- Улучшенный мониторинг с трейсингом запросов
- Метрики и алерты для отслеживания состояния системы
- Централизованный сбор и анализ логов
- Более сложную систему ролей и прав доступа
- Интеграцию с другими сервисами экосистемы

Это позволит создать более надёжную, безопасную и масштабируемую систему, способную обслуживать растущее количество пользователей и запросов. 