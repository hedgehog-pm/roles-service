# Архитектура и структура

Архитектура опирается на правила и скиллы репозитория и на
[service-tracker](https://github.com/hedgehog-pm/service-tracker/tree/1b9a0022cb6734a86b76be0ddffc268b8263515d).
Сравнение выполнено с ревизией `1b9a0022cb6734a86b76be0ddffc268b8263515d`.
Предметные правила roles-service сохраняются при изменении структуры.

## Основания решений

| Решение | Источник |
| --- | --- |
| Сборка зависимостей и жизненный цикл в cmd/service | `.agents/rules/architecture.md`, `golang-design-patterns`, [точка входа tracker](https://github.com/hedgehog-pm/service-tracker/blob/1b9a0022cb6734a86b76be0ddffc268b8263515d/cmd/service/main.go) |
| Простые операции группируются по сущности | `.agents/rules/usecase.md`, [task.Usecase](https://github.com/hedgehog-pm/service-tracker/blob/1b9a0022cb6734a86b76be0ddffc268b8263515d/internal/usecase/task/usecase.go) |
| Сценарий с особыми зависимостями имеет New/Do и собственные Params | `.agents/rules/usecase.md`, [listtasks](https://github.com/hedgehog-pm/service-tracker/tree/1b9a0022cb6734a86b76be0ddffc268b8263515d/internal/usecase/listtasks) |
| Контракт зависимости принадлежит потребителю | `golang-design-patterns`, [интерфейсы listtasks](https://github.com/hedgehog-pm/service-tracker/blob/1b9a0022cb6734a86b76be0ddffc268b8263515d/internal/usecase/listtasks/interface.go) |
| Агрегаты, отдельный интерфейс для составного чтения, транзакции usecase | `.agents/rules/repository.md` |
| Фильтры и страницы находятся вне домена | `.agents/rules/repository.md`, [Params и Page listboardcolumntasks](https://github.com/hedgehog-pm/service-tracker/blob/1b9a0022cb6734a86b76be0ddffc268b8263515d/internal/usecase/listboardcolumntasks/usecase.go) |
| Перенос сохраняет порядок проверок, ошибки и побочные эффекты | `.agents/rules/coding.md`, `golang-code-style`, `golang-error-handling` |

Правила находятся в `.agents/rules`, скиллы — в `.agents/skills/<имя>/SKILL.md`.
Общие типы пагинации вынесены в небольшой пакет, потому что три списка уже
используют одинаковые границы, значения по умолчанию и формат страницы.
Он не выполняет запросы и не образует дополнительный слой хранения.

## Структура

```text
cmd/
  service/                   сборка сервиса и завершение процесса
  migrate/                   запуск миграций
internal/
  api/                       HTTP, аутентификация, DTO и ошибки транспорта
  domain/                    Role, RoleTemplate, RoleAssignment, права и scope
  pagination/                параметры страниц, сортировка и Page[T]
  usecase/
    role/                    Get, Delete
    roletemplate/            Get, Page, Delete, CheckSystem
    roleassignment/          List, Assign, Revoke
    permission/              каталог прав
    createrole/              создание роли с проверкой каталога
    updaterole/              правила изменения роли и её версии
    createtemplate/          создание шаблона с проверкой каталога
    updatetemplate/          правила изменения шаблона и его версии
    listroles/               список с проверкой фильтра по каталогу
    listroleassignments/     составное чтение роли и назначения
  repository/
    role/                    роль вместе с правами и outbox
    roletemplate/            шаблон вместе с правами
    roleassignment/          назначения и outbox; отдельный Reader
  storage/postgres/          пул, транзакции, executor, преобразование ошибок
  generated/
    api/                     код из OpenAPI
    bob/                     SQL-код по группам role/roletemplate/roleassignment
api/server/roles/            OpenAPI и конфигурация генерации
docs/                        документация backend и frontend
```

В каждом usecase зависимости объявлены в `interface.go`. Конструкторы
получают обязательные зависимости явно. В `cmd/service/main.go` видна вся
сборка: пул → Transactor → repository/Reader → usecase → handlers → сервер.

## Границы и зависимости

| Пакет | Ответственность и разрешённые зависимости |
| --- | --- |
| domain | Состояние агрегатов и предметные инварианты; только стандартная библиотека |
| pagination | Общий контракт списка и прежние ошибки валидации через domain |
| usecase | Последовательность действий, интерфейсы внешних эффектов и границы транзакций; domain и pagination |
| api | Преобразование запросов и ответов, вызов usecase; без SQL и repository |
| repository | SQL, восстановление агрегатов, сохранение полного состояния, outbox; domain, Bob и storage/postgres |
| storage/postgres | Работа с PostgreSQL и перевод технических ошибок; без зависимостей от usecase и API |
| cmd/service | Создание конкретных реализаций и управление временем их жизни |

Репозитории также используют объявленные потребителем типы фильтров и результатов
usecase. Это зависимость реализации от контракта чтения: usecase не импортирует
repository, PostgreSQL или generated. Дополнительных копий Role/RoleTemplate нет.

Роль и её права сохраняются целиком одним repository. Шаблон и его права —
другим. Назначение ссылается на роль по ID и имеет собственный repository.
Статический PermissionCatalog не требует repository.

`listroleassignments.reader` возвращает собственную проекцию `Assignment`:
она содержит имя и scope роли, которых нет в агрегате RoleAssignment.
Её PostgreSQL-реализация — `roleassignment.Reader` в `reader.go`.
Параметры фильтра находятся рядом с потребителем. SQL чтения остаётся в
`internal/repository/roleassignment/queries.sql`.

Удаления входят в entity-usecase, поскольку используют тот же repository и
Transactor. Assign/Revoke/List также работают с одним repository назначений.
Создание и обновление ролей/шаблонов используют дополнительные зависимости
и правила каталога. `listroles` имеет дополнительную зависимость от каталога
для проверки permission_code. Список шаблонов остаётся методом roletemplate:
отдельного набора зависимостей для него нет.

Проверка принадлежности шаблона системе выполняется в `roletemplate.CheckSystem`.
HTTP вызывает её в прежнем порядке перед чтением, изменением или удалением.
Чужая система по-прежнему получает `not_found`.

## Сохранённое поведение

- Доменные конструкторы, изменение имени/прав, каталог и ошибки не изменены.
- OpenAPI, SQL-запросы, миграции и сгенерированные контракты сохранены.
- Значения по умолчанию, пределы и порядок валидации пагинации сохранены.
- Создание, обновление, удаление, назначение и отзыв выполняются в `InTx`.
- Страницы используют `InReadTx` с RepeatableRead и ReadOnly: Count и Items
  читаются в одном снимке. Это сохраняет имеющийся контракт roles-service.
- Проверка ожидаемой версии и блокировки строк сохранены. Обновление двух
  полей увеличивает версию один раз; запрос без изменений её не увеличивает.
- Повторное назначение возвращает исходные автора и время. Повторный отзыв
  существующей роли успешен и не создаёт лишнее событие.
- Запись outbox атомарна с изменением роли или назначения. Шаблоны не создают
  собственных событий. При ошибке outbox изменения откатываются.
- Transactor передаёт транзакцию через контекст; repository выбирает executor
  из этого контекста и не открывает свою транзакцию.

Подробности HTTP и интеграции находятся в [руководстве frontend](guides/frontend.md),
ограничения аутентификации и внешней авторизации — в [ограничениях](limitations.md).

## Инструменты и проверка

`.golangci.yml`, `AGENTS.md`, четыре файла `.agents/rules` и 51 файл
`.agents/skills` совпадают с указанной ревизией service-tracker.
В обоих проектах golangci-lint v2.15.0; цели fmt/lint/fix вызывают одинаковые
команды. Дополнительных правил go-arch-lint/Ruleguard нет.
MkDocs и Swagger-плагин совпадают с service-tracker. Material 9.8.0 используется
на Python 3.11 и новее; для установленного в Windows Python 3.9 закреплена
совместимая версия 9.7.7. Hook копирования OpenAPI совпадает с эталоном.

Проверки проекта: `make fmt`, `make lint`, `go build ./cmd/...`, `go test ./...`,
`mkdocs build --strict`. Если make отсутствует, выполняются команды его целей.
Проверка импортов дополняет линтер: domain не зависит от других внутренних
пакетов, usecase не импортирует транспорт и БД, API не импортирует repository.

Тестовые файлы отсутствуют по решению пользователя. `go test ./...` в этом
состоянии проверяет компилируемость пакетов и сообщает `no test files`.
Проверки HTTP и сохранения данных проводятся на отдельной временной PostgreSQL;
они не доказывают корректность всех конкурентных сценариев.
