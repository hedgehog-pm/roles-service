# Интеграция фронтенда

## Маршруты и обязательные параметры

| Метод и маршрут | Параметры или тело | Успех |
| --- | --- | --- |
| `GET /permissions` | query `scope_type` | `200`, массив Permission |
| `GET /roles` | query `scope_type`, `scope_id` | `200`, RolePage |
| `POST /roles` | JSON `name`, `scope`, `permission_codes` | `201`, Role |
| `GET /roles/{id}` | path `id` | `200`, Role |
| `PUT /roles/{id}` | JSON `name`, `permission_codes`, `expected_version` | `200`, Role |
| `DELETE /roles/{id}` | query `expected_version` | `204` |
| `GET /role-templates` | query `system_id`, `target_scope_type` | `200`, TemplatePage |
| `POST /role-templates` | JSON `name`, `system_id`, `target_scope_type`, `permission_codes` | `201`, RoleTemplate |
| `GET /role-templates/{id}` | query `system_id` | `200`, RoleTemplate |
| `PUT /role-templates/{id}` | query `system_id`; JSON `name`, `permission_codes`, `expected_version` | `200`, RoleTemplate |
| `DELETE /role-templates/{id}` | query `system_id`, `expected_version` | `204` |
| `GET /roles/{role_id}/assignments` | path `role_id` | `200`, AssignmentPage |
| `PUT /roles/{role_id}/assignments/{user_id}` | path IDs, без тела | `201` или `200`, RoleAssignment |
| `DELETE /roles/{role_id}/assignments/{user_id}` | path IDs | `204` |
| `GET /participants/{user_id}/role-assignments` | query `scope_type`, `scope_id` | `200`, UserRoleAssignmentPage |
| `GET /role-assignments` | query `scope_type`, `scope_id` | `200`, UserRoleAssignmentPage |

ID — целые числа от 1 до 2147483647. `scope_type` и `target_scope_type`: `system`, `workspace`, `project`. `scope_id` обозначает ресурс соответствующего типа. У шаблона `system_id` — контекст хранения шаблона, а `target_scope_type` — тип будущей роли. Для GET/PUT/DELETE шаблона чужой `system_id` даёт `404`.

## Формы ролей и шаблонов

1. Получите каталог для выбранного scope.
2. Покажите имя и выбор нескольких прав. Требуется хотя бы одно право.
3. Отправьте только коды выбранных прав в `permission_codes`.
4. После сохранения замените локальный объект ответом сервера, включая `version`.

Имя обрезается по краям и должно содержать от 1 до 200 Unicode-символов. Строка из пробелов недопустима. Коды должны быть уникальны и принадлежать scope. Не отправляйте дополнительные поля: схемы запросов их запрещают. `created_by`, `assigned_by`, даты и ID создаёт backend.

При редактировании получите коды через `role.permissions.map(p => p.code)`. Ответ содержит `permissions`, запрос — `permission_codes`. `PUT` требует оба изменяемых поля даже при изменении только имени. Scope роли, `system_id` и целевой scope шаблона через PUT не меняются.

Даты передаются строками RFC 3339 с часовым поясом. Для отображения преобразуйте их в локальное время пользователя, исходное значение сохраняйте без округления.

Для создания роли по шаблону скопируйте в форму имя и коды, выберите конкретный scope и выполните обычный `POST /roles`. Изменение шаблона не меняет ранее созданные роли.

## Версии и конфликты

Для обновления передавайте `expected_version` в JSON, для удаления — в query. Берите значение из последнего ответа API, не вычисляйте версию на клиенте. Изменение имени и прав в одном запросе увеличивает итоговую версию один раз. Запрос без изменений сохраняет её.

На `409 version_conflict` сохраните черновик формы, повторно загрузите объект и предложите сравнить изменения. Не повторяйте запись автоматически с новой версией: это может затереть чужую работу. `404` при записи означает, что объект отсутствует либо не соответствует контексту шаблона.

## Списки

Все списки, кроме каталога прав, имеют форму:

```json
{"items": [], "total": 0, "limit": 20, "offset": 0}
```

`limit`: 1–100, по умолчанию 20. `offset`: 0–2147483647, по умолчанию 0. `total` учитывает фильтры до пагинации. Пустая страница при большом offset может иметь положительный total. При смене фильтров сбрасывайте offset в 0. Count и items одного ответа читаются из общего снимка; разные HTTP-запросы могут видеть изменения данных.

Дополнительные параметры ролей и шаблонов:

| Параметр | Поведение |
| --- | --- |
| `search` | Подстрока имени без учёта регистра; `%` и `_` — обычные символы |
| `sort` | `name`, `created_at`, `updated_at`; по умолчанию `created_at` |
| `direction` | `asc`, `desc`; по умолчанию `desc` |
| `permission_code` | Только роли: содержит заданный код права |

При равных значениях сортировки используется ID в том же направлении. Назначения сортируются по `role_id`, затем `user_id`; параметров search/sort у них нет.

Список назначений роли возвращает `role_id`, `user_id`, `assigned_by`, `assigned_at`. Списки пользователя и scope дополнительно возвращают `role_name` и `scope`. Имена пользователей сервис не отдаёт — их получает UI из сервиса пользователей. Пустой список назначений роли сам по себе не доказывает существование роли.

## Ошибки

```json
{
  "code": "invalid_input",
  "message": "invalid input"
}
```

| HTTP | code | Действие UI |
| --- | --- | --- |
| 400 | `invalid_input` | Показать ошибку поля или общую ошибку формы |
| 400 | `unknown_permission_code` | Обновить каталог и проверить выбранные коды |
| 400 | `permission_scope_mismatch` | Проверить scope и выбранные права |
| 400 | `empty_permissions` | Потребовать хотя бы одно право |
| 404 | `not_found` | Убрать устаревший объект либо проверить system_id |
| 409 | `version_conflict` | Загрузить актуальную версию и согласовать черновик |
| 500 | `internal_error` | Показать общую ошибку, технические детали не доступны |

Ошибки схемы OpenAPI возвращают `invalid_input` без `details`, в том числе при пустом `permission_codes`. Доменные ошибки могут содержать `details`. Логику стройте по HTTP и `code`, а не по английскому `message`.

## Минимальный fetch-клиент

```typescript
type ApiProblem = {
  code: string;
  message: string;
  details?: { field: string; code: string; message: string }[];
};

export class ApiError extends Error {
  constructor(public status: number, public problem: ApiProblem) {
    super(problem.message);
  }
}

export async function request<T>(base: string, path: string, init: RequestInit = {}): Promise<T | undefined> {
  const headers = new Headers(init.headers);
  if (init.body !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  const response = await fetch(base + path, { ...init, headers });
  if (response.status === 204) return undefined;
  const body: unknown = await response.json();
  if (!response.ok) throw new ApiError(response.status, body as ApiProblem);
  return body as T;
}
```

Пример рассчитан на JSON-ответы самого сервиса. Ошибки сети, отмену `AbortSignal` и ответы прокси не в JSON обрабатывайте отдельно. `as T` не выполняет runtime-валидацию; модели удобно получать генератором из OpenAPI. Для query используйте `URLSearchParams`. Не добавляйте заголовок actor ID: такой механизм не поддерживается.
