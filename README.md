# go-common-config

Универсальный пакет конфигурации для Go-сервисов в едином формате.

## Что внутри

- загрузка конфигурации из переменных окружения (с поддержкой `.env` для локалки)
- нормализация значений
- fail-fast валидация
- генерация PostgreSQL DSN

## Подключение из другого сервиса

```bash
go get github.com/khrll/go-common-config@v1.0.0
```

```go
import config "github.com/khrll/go-common-config"

cfg, err := config.LoadConfig()
if err != nil {
    // обработка ошибки
}
```

## Локальная разработка до публикации

В сервисе можно использовать `replace` в `go.mod`:

```go
require github.com/khrll/go-common-config v0.0.0
replace github.com/khrll/go-common-config => ./packages/go-common-config
```

После публикации в Git удалите `replace` и укажите версию тега.
