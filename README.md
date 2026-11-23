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

- [Для фронтендеров](docs/index.md) — быстрый старт, формы, списки, версии и ошибки.
- [Интеграция фронтенда](docs/guides/frontend.md) — маршруты и пример fetch-клиента.
- [OpenAPI-контракт](api/server/roles/openapi.yaml) — маршруты, запросы и ответы.
- [Описание модуля](docs/role-module.md) — сценарии, хранение, транзакции,
  валидация, ошибки и генерация кода.

OpenAPI-файл можно импортировать в Postman

Документация собирается MkDocs Material со Swagger. Команды одинаковы в
Windows, Linux и macOS:

```sh
make docs
```

После запуска откройте `http://127.0.0.1:8000`. Строгая сборка без запуска
локального сервера выполняется командой `make docs-build`. Make сам создаёт
`.venv-docs` и устанавливает зависимости.
Python 3.9 использует совместимый Material 9.7.7, Python 3.11 и новее — 9.8.0.

Публикация сайта автоматически не выполняется. OpenAPI копируется в сайт
хуком `scripts/docs.py` из исходного контракта.

### Проверки в PowerShell без make

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
