# Giper SaaS Common Libraries

Общие переиспользуемые библиотеки для всех микросервисов платформы Giper SaaS.

## Структура

```
common/
├── pkg/
│   └── httputil/          # HTTP утилиты для трассировки запросов
│       └── context.go     # Request ID propagation
├── go.mod
└── README.md
```

## Пакеты

### pkg/httputil

HTTP утилиты для distributed request tracing с использованием `X-Request-ID` и `X-Correlation-ID` заголовков.

**Основные функции:**

- `ContextFromGin(c)` - Извлекает request_id из gin.Context и создает context.Context
- `PropagateRequestIDFromContext(ctx, req)` - Добавляет заголовки к исходящим HTTP-запросам
- `GetRequestID(c)` - Извлекает request_id из gin.Context
- `GetRequestIDFromContext(ctx)` - Извлекает request_id из context.Context
- `ContextWithRequestID(ctx, id)` - Создает context с request_id

**Константы:**

- `HeaderRequestID` - "X-Request-ID"
- `HeaderCorrelationID` - "X-Correlation-ID"
- `RequestIDKey` - "request_id" (для gin.Context)

## Использование

### Установка

Добавьте в `go.mod` вашего сервиса:

```go
require (
    gitlab.vertical-tech.ru/gipersass/common v0.1.0
)
```

Или используйте replace directive для локальной разработки:

```go
replace gitlab.vertical-tech.ru/gipersass/common => ../common
```

### Обновление зависимостей

```bash
# В корне вашего сервиса (erp-service, user-service, etc.)
go get gitlab.vertical-tech.ru/gipersass/common@latest
# или
go mod tidy
```

### Пример: Handler с request tracing

```go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gitlab.vertical-tech.ru/gipersass/common/pkg/httputil"
)

func (h *Handler) CreateOrder(c *gin.Context) {
    var req OrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Получаем контекст с request_id
    ctx := httputil.ContextFromGin(c)

    // Передаем в service layer
    order, err := h.service.CreateOrder(ctx, req)
    if err != nil {
        h.log.Error("failed to create order", "error", err,
            "request_id", httputil.GetRequestIDFromContext(ctx))
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(200, order)
}
```

### Пример: HTTP Client с request tracing

```go
package client

import (
    "context"
    "net/http"

    "gitlab.vertical-tech.ru/gipersass/common/pkg/httputil"
)

func (c *Client) CallExternalAPI(ctx context.Context, url string) error {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return err
    }

    // Пробрасываем request_id в заголовках
    httputil.PropagateRequestIDFromContext(ctx, req)

    c.log.Info("calling external API",
        "url", url,
        "request_id", httputil.GetRequestIDFromContext(ctx),
    )

    resp, err := c.httpClient.Do(req)
    // ...

    return nil
}
```

### Пример: Middleware

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gitlab.vertical-tech.ru/gipersass/common/pkg/httputil"
)

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Проверяем заголовки
        requestID := c.GetHeader(httputil.HeaderRequestID)
        if requestID == "" {
            requestID = c.GetHeader(httputil.HeaderCorrelationID)
        }
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // Сохраняем в контекст
        c.Set(httputil.RequestIDKey, requestID)

        // Возвращаем в ответе
        c.Writer.Header().Set(httputil.HeaderRequestID, requestID)
        c.Writer.Header().Set(httputil.HeaderCorrelationID, requestID)

        c.Next()
    }
}
```

## Версионирование

Следуем [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** версия при несовместимых изменениях API
- **MINOR** версия при добавлении функциональности с обратной совместимостью
- **PATCH** версия при исправлении багов с обратной совместимостью

## Changelog

### v0.1.0 (2026-02-18)

- ✨ Добавлен пакет `httputil` для distributed request tracing
- ✨ Поддержка `X-Request-ID` и `X-Correlation-ID` заголовков
- 📝 Добавлена документация и примеры использования

## Сервисы, использующие common

- ✅ `erp-service` - ERP интеграция (1С, товары, заказы)
- ⏳ `user-service` - Управление пользователями
- ⏳ `ozon-service` - Ozon маркетплейс интеграция

## Разработка

### Добавление нового пакета

1. Создайте директорию в `pkg/`:
```bash
mkdir -p pkg/newpackage
```

2. Добавьте код и тесты:
```bash
touch pkg/newpackage/newpackage.go
touch pkg/newpackage/newpackage_test.go
```

3. Обновите README с описанием нового пакета

4. Создайте git tag для новой версии:
```bash
git tag v0.2.0
git push origin v0.2.0
```

### Тестирование

```bash
# Запустить все тесты
go test ./...

# С покрытием
go test -cover ./...

# Линтер
golangci-lint run
```

### Локальная разработка

Для тестирования изменений до публикации используйте replace directive:

```go
// В go.mod вашего сервиса
replace gitlab.vertical-tech.ru/gipersass/common => ../common
```

После тестирования создайте git tag и обновите сервисы:

```bash
# В common/
git tag v0.2.0
git push origin v0.2.0

# В каждом сервисе
go get gitlab.vertical-tech.ru/gipersass/common@v0.2.0
```

## Лицензия

Proprietary - Giper SaaS Platform

## Контакты

- Команда: Giper SaaS Development Team
- Репозиторий: gitlab.vertical-tech.ru/gipersass/common