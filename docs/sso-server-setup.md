# Настройка и запуск gRPC-сервера в SSO

В проекте SSO реализована система настройки и запуска gRPC-сервера, которая обеспечивает надежную обработку запросов аутентификации. Давайте рассмотрим основные компоненты этой системы.

## Регистрация сервиса аутентификации

Функция `authgrpc.Register` играет ключевую роль в настройке сервера. Она связывает реализацию сервиса аутентификации с gRPC-сервером:

```go
func Register(gRPCServer *grpc.Server, auth Auth) {
    ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}
```

Это позволяет серверу обрабатывать входящие RPC-запросы, связанные с аутентификацией. Методы сервиса (Login, Register, IsAdmin) становятся доступными для клиентов через gRPC.

## Запуск сервера

В проекте реализованы две функции для запуска сервера:

```go
// MustRun запускает gRPC-сервер и вызывает панику при ошибке
func (a *App) MustRun() {
    if err := a.Run(); err != nil {
        panic(err)
    }
}

// Run запускает gRPC-сервер
func (a *App) Run() error {
    const op = "grpcapp.Run"

    // Создание TCP-листенера
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
    if err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

    // Запуск обработчика gRPC-сообщений
    if err := a.gRPCServer.Serve(l); err != nil {
        return fmt.Errorf("%s: %w", op, err)
    }

    return nil
}
```

## Структура приложения

Основное приложение SSO организовано следующим образом:

```go
type App struct {
    GRPCServer *grpcapp.App
}

func New(
    log *slog.Logger,
    grpcPort int,
    storagePath string,
    tokenTTL time.Duration,
) *App {
    // Инициализация хранилища
    storage, err := sqlite.New(storagePath)
    if err != nil {
        panic(err)
    }

    // Создание сервиса аутентификации
    authService := auth.New(log, storage, storage, storage, tokenTTL)

    // Создание gRPC-приложения
    grpcApp := grpcapp.New(log, authService, grpcPort)

    return &App{
        GRPCServer: grpcApp,
    }
}
```

## Точка входа

Запуск приложения происходит в `cmd/sso/main.go`:

```go
func main() {
    // Загрузка конфигурации
    cfg := config.MustLoad()

    // Настройка логгера
    log := setupLogger(cfg.Env)

    // Создание приложения
    application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

    // Запуск сервера
    application.GRPCServer.MustRun()
}
```

## Особенности реализации

1. **Обработка ошибок**:
   - Использование `MustRun()` для немедленной остановки при критических ошибках
   - Детальное логирование ошибок с контекстом
   - Graceful shutdown при корректном завершении

2. **Конфигурация**:
   - Гибкая настройка через конфигурационный файл
   - Поддержка переменных окружения
   - Настраиваемые таймауты и порты

3. **Логирование**:
   - Структурированные логи с контекстом
   - Информация о запуске и остановке сервера
   - Отслеживание состояния соединений

4. **Безопасность**:
   - Защита от паники через recovery-интерцептор
   - Безопасное хранение токенов
   - Валидация входных данных

## Рекомендации по улучшению

1. **Мониторинг**:
   - Добавить метрики состояния сервера
   - Реализовать health check endpoints
   - Настроить алерты при проблемах

2. **Масштабируемость**:
   - Подготовить к горизонтальному масштабированию
   - Реализовать балансировку нагрузки
   - Добавить кэширование

3. **Безопасность**:
   - Усилить валидацию входных данных
   - Добавить rate limiting
   - Реализовать механизм блокировки при подозрительной активности 