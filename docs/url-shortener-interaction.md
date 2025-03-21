# Схема взаимодействия компонентов URL Shortener

Схема взаимодействия компонентов URL Shortener представляет собой мВногоуровневую архитектуру, где каждый компонент имеет четко определенную роль в обработке запросов. Процесс начинается с HTTP-клиента, который отправляет запросы на создание коротких ссылок или их получение. Запросы проходят через Chi Router, который маршрутизирует их к соответствующим обработчикам, предварительно обрабатывая через цепочку middleware.

 цепочке middleware запрос последовательно проходит через RequestID Middleware, добавляющий уникальный идентификатор к каждому запросу, Logger Middleware для записи деталей запроса, Recoverer Middleware для обработки паник и URLFormat Middleware для нормализации URL-путей. После обработки middleware запрос попадает к одному из двух обработчиков: Save Handler для создания новых коротких ссылок или Redirect Handler для получения оригинального URL по алиасу.

Save Handler выполняет валидацию входных данных через Validator, проверяя корректность URL и наличие обязательных полей. При успешной валидации запрос передается в URL Service, который реализует бизнес-логику и взаимодействует с хранилищем данных. URL Service также проверяет права доступа через SSO Client, который общается с сервисом SSO по протоколу gRPC.

Взаимодействие с хранилищем данных происходит через SQLite, где сохраняются маппинги между оригинальными URL и их алиасами. Хранилище обеспечивает уникальность алиасов и предоставляет методы для сохранения и получения URL. После обработки запроса формируется соответствующий ответ, который может быть успешным с созданным алиасом, сообщением об ошибке, перенаправлением на оригинальный URL или 404-ошибкой при отсутствии запрашиваемого URL.

```mermaid
graph TD
    Client[HTTP Client] -->|POST /url, GET /{alias}| Router[Chi Router]
    
    subgraph Middleware
        Router -->|RequestID| MW1[RequestID Middleware]
        MW1 -->|Logger| MW2[Logger Middleware]
        MW2 -->|Recoverer| MW3[Recoverer Middleware]
        MW3 -->|URLFormat| MW4[URLFormat Middleware]
    end
    
    subgraph Handlers
        MW4 -->|Save URL| H1[Save Handler]
        MW4 -->|Redirect| H2[Redirect Handler]
        
        H1 -->|Validate Request| V1[Validator]
        V1 -->|Valid| S1[URL Service]
        V1 -->|Invalid| E1[Error Response]
        
        H2 -->|Get URL| S2[URL Service]
    end
    
    subgraph Services
        S1 -->|Save| DB[(SQLite Storage)]
        S2 -->|Get| DB
    end
    
    subgraph SSO Integration
        S1 -->|Check Admin| SSO[SSO Client]
        SSO -->|gRPC| SSOServer[SSO Service]
    end
    
    subgraph Response
        S1 -->|Success| R1[Success Response]
        S1 -->|Error| R2[Error Response]
        S2 -->|Found| R3[Redirect Response]
        S2 -->|NotFound| R4[404 Response]
    end
    
    R1 --> Client
    R2 --> Client
    R3 --> Client
    R4 --> Client
    E1 --> Client
```

## Описание компонентов

### HTTP Client
- Отправляет запросы на создание коротких ссылок (POST /url)
- Отправляет запросы на перенаправление (GET /{alias})
- Получает JSON-ответы с результатами операций

### Chi Router
- Маршрутизирует запросы к соответствующим обработчикам
- Применяет цепочку middleware для каждого запроса

### Middleware
1. **RequestID Middleware**
   - Добавляет уникальный ID к каждому запросу
   - Используется для отслеживания запросов в логах

2. **Logger Middleware**
   - Логирует детали каждого HTTP-запроса
   - Записывает метод, URL, статус ответа

3. **Recoverer Middleware**
   - Обрабатывает паники в обработчиках
   - Возвращает 500 ошибку при панике

4. **URLFormat Middleware**
   - Нормализует URL-пути
   - Обрабатывает trailing slashes

### Handlers
1. **Save Handler**
   - Обрабатывает POST /url запросы
   - Валидирует входные данные
   - Сохраняет URL в хранилище

2. **Redirect Handler**
   - Обрабатывает GET /{alias} запросы
   - Получает оригинальный URL
   - Выполняет перенаправление

### Services
1. **URL Service**
   - Реализует бизнес-логику
   - Взаимодействует с хранилищем
   - Проверяет права через SSO

### Storage
- SQLite база данных
- Хранит маппинги URL и алиасов
- Обеспечивает уникальность алиасов

### SSO Integration
- Проверяет права доступа через gRPC
- Интегрируется с сервисом SSO
- Обрабатывает ошибки авторизации

### Response
1. **Success Response**
   ```json
   {
     "status": "OK",
     "alias": "generated_alias"
   }
   ```

2. **Error Response**
   ```json
   {
     "status": "Error",
     "error": "error message"
   }
   ```

3. **Redirect Response**
   - HTTP 302 Found
   - Location header с оригинальным URL

4. **404 Response**
   ```json
   {
     "status": "Error",
     "error": "URL not found"
   }
   ```

## Обработка ошибок

На каждом уровне происходит обработка ошибок:

```go
// Валидация
func ValidationError(errs validator.ValidationErrors) Response {
    var errMsgs []string
    for _, err := range errs {
        switch err.ActualTag() {
        case "required":
            errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
        case "url":
            errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
        default:
            errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
        }
    }
    return Response{Status: StatusError, Error: strings.Join(errMsgs, ", ")}
}

// Обработка ошибок хранилища
if err != nil {
    log.Error("failed to add url", sl.Err(err))
    render.JSON(w, r, resp.Error("failed to add url"))
    return
}
```

## Логирование

Логирование происходит на каждом этапе обработки запроса:

```go
log.Info("request body decoded", slog.Any("req", req))
log.Info("url added", slog.Int64("id", id))
log.Error("failed to decode request body", sl.Err(err))
```

## Интеграция с SSO

При необходимости проверки прав доступа:

```go
// Проверка прав через SSO клиент
isAdmin, err := ssoClient.IsAdmin(ctx, token)
if err != nil {
    log.Error("failed to check admin rights", sl.Err(err))
    render.JSON(w, r, resp.Error("failed to check admin rights"))
    return
}
``` 