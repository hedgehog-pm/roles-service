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

go run ./cmd/migrate
go run ./cmd/service
```

Пока сервис аутентификации не готов, подключена middleware `DevelopmentActor`.
Все операции выполняются от пользователя с ID `1`, Bearer-токен не требуется
и не проверяется. Переменная `DEV_ACTOR_ID` не используется.
Middleware `Authentication` подготовлена для подключения проверки токена,
но сейчас не включена. Миграции выполняйте отдельно перед запуском сервиса.

## Документация

Сервис находится в разработке. Документацию подготовим после завершения
разработки; существующие материалы в `docs/` могут быть неактуальны.

Текущий [OpenAPI-контракт](api/server/roles/openapi.yaml) можно импортировать
в Postman.

## Проверки в PowerShell без make

Запускайте из корня репозитория. Команды соответствуют целям Makefile:

~~~powershell
go tool golangci-lint fmt
go tool golangci-lint config verify
go tool golangci-lint run
go build ./cmd/...
go test ./...
~~~

Проверяйте код завершения каждой команды; при ошибке остановитесь.

## Архитектура

Прикладные сценарии находятся в internal/usecase, репозитории агрегатов —
в internal/repository/{role,roletemplate,roleassignment}. Usecase задаёт транзакции,
а internal/storage/postgres передаёт их через context. SQL лежит рядом с
репозиториями; Bob генерирует отдельные пакеты в internal/generated/bob.
make generate запускает генерацию SQL и API; SQL требует мигрированную БД.

Линтер и его конфигурация совпадают с service-tracker: golangci-lint v2.15.0.
