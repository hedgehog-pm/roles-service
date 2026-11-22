# roles-service

Сервис управления ролями, шаблонами ролей, назначениями пользователей и
каталогом прав для system/workspace/project.

## Локальный запуск

Нужны Go и запущенный Docker Desktop. Выполните из корня проекта в PowerShell:

```powershell
$env:ROLES_DATABASE_PASSWORD = "roles_local"
docker compose up -d roles-db

$env:DATABASE_URL = "postgres://roles_service:roles_local@127.0.0.1:55432/roles_service?sslmode=disable"
$env:HTTP_ADDR = "127.0.0.1:8080"
$env:DEV_ACTOR_ID = "1"

go run ./cmd/migrate
go run ./cmd/service
```

`DEV_ACTOR_ID` используется только для локального запуска без подключённой
аутентификации. Миграции необходимо выполнять отдельно перед запуском сервиса.

## Документация

- [OpenAPI-контракт](api/server/roles/openapi.yaml) — маршруты, запросы и ответы.
- [Описание модуля](docs/role-module.md) — сценарии, хранение, транзакции,
  валидация, ошибки и генерация кода.

OpenAPI-файл можно импортировать в Postman

### Проверки в PowerShell без make

Запускайте из корня репозитория. Команды соответствуют целям Makefile:

~~~powershell
go tool golangci-lint fmt
go tool golangci-lint config verify
go tool golangci-lint run
go run github.com/fe3dback/go-arch-lint@v1.19.0 check
go test ./...
~~~

Проверяйте код завершения каждой команды; при ошибке остановитесь.