# Логика работы приложения SSO + URL Shortener

```mermaid
flowchart TD
    Start[Начало] --> Client[Клиент]
    
    subgraph Клиентская часть
        Client -->|1. Отправка запроса| Request[HTTP POST /api/urls]
        Request -->|2. Добавление токена| Token[Bearer JWT Token]
    end
    
    subgraph URL Shortener
        Token -->|3. Проверка токена| ValidateToken{Валидный токен?}
        ValidateToken -->|Нет| Error401[401 Unauthorized]
        ValidateToken -->|Да| ExtractUserID[Извлечение userID]
        ExtractUserID -->|4. Проверка прав| CheckAdmin{Администратор?}
        CheckAdmin -->|Нет| Error403[403 Forbidden]
        CheckAdmin -->|Да| ProcessURL[Обработка URL]
    end
    
    subgraph SSO Service
        CheckAdmin -->|5. gRPC запрос| IsAdmin[IsAdmin RPC]
        IsAdmin -->|6. Проверка БД| DBQuery[SQLite Query]
        DBQuery -->|7. Результат| AdminCheck{is_admin?}
        AdminCheck -->|Нет| ReturnFalse[Возврат false]
        AdminCheck -->|Да| ReturnTrue[Возврат true]
    end
    
    subgraph Обработка результата
        ReturnFalse -->|8. Отказ в доступе| Error403
        ReturnTrue -->|9. Разрешение доступа| ProcessURL
        ProcessURL -->|10. Сохранение URL| SaveURL[Сохранение в БД]
        SaveURL -->|11. Ответ клиенту| Response[HTTP 200 OK]
        Error401 -->|12. Ошибка| ErrorResponse[HTTP Error Response]
        Error403 -->|12. Ошибка| ErrorResponse
    end
    
    Response --> End[Конец]
    ErrorResponse --> End
    
    style Start fill:#f9f,stroke:#333,stroke-width:2px
    style End fill:#f9f,stroke:#333,stroke-width:2px
    style ValidateToken fill:#fbb,stroke:#333,stroke-width:2px
    style CheckAdmin fill:#fbb,stroke:#333,stroke-width:2px
    style AdminCheck fill:#fbb,stroke:#333,stroke-width:2px
    style Error401 fill:#f99,stroke:#333,stroke-width:2px
    style Error403 fill:#f99,stroke:#333,stroke-width:2px
    style Response fill:#9f9,stroke:#333,stroke-width:2px
```

## Описание логики работы

1. **Инициация запроса**
   - Клиент формирует HTTP POST запрос
   - Добавляет JWT токен в заголовок Authorization

2. **Проверка токена**
   - URL Shortener проверяет наличие токена
   - Валидирует JWT токен
   - При ошибке возвращает 401 Unauthorized

3. **Извлечение данных**
   - Из токена извлекается userID
   - Подготавливается запрос к SSO

4. **Проверка прав**
   - URL Shortener отправляет gRPC запрос к SSO
   - SSO проверяет права в базе данных
   - При отсутствии прав возвращает 403 Forbidden

5. **Обработка URL**
   - При успешной проверке прав обрабатывается URL
   - URL сохраняется в базе данных
   - Возвращается успешный ответ клиенту

6. **Обработка ошибок**
   - Все ошибки логируются
   - Клиент получает понятное сообщение об ошибке
   - Сохраняется консистентность данных 