# Мини-схема архитектуры SSO

```mermaid
%%{init: {'theme': 'base', 'themeVariables': { 'fontSize': '14px', 'fontFamily': 'arial'}}}%%
graph TB
    subgraph Client["Клиент"]
        C[gRPC Клиент]
    end
    
    subgraph Server["Сервер"]
        A[Auth Service]
        J[JWT]
        S[(SQLite)]
    end
    
    C -->|Login/Register/IsAdmin| A
    A -->|Генерация| J
    A -->|CRUD| S
    
    classDef client fill:#f9f,stroke:#333,stroke-width:2px;
    classDef server fill:#bbf,stroke:#333,stroke-width:2px;
    
    class Client client;
    class Server server;
```

## Краткое описание

Система SSO построена на трех основных компонентах: клиент, сервер аутентификации и хранилище данных. Клиент взаимодействует с сервером через gRPC API, отправляя запросы на аутентификацию, регистрацию и проверку прав. Сервер обрабатывает эти запросы, используя JWT для генерации токенов и SQLite для хранения данных. Вся система построена с учетом принципов чистой архитектуры и безопасности. 