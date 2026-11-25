-- GetRole
SELECT r.*, ARRAY(
    SELECT p.permission_code FROM role_permissions p
    WHERE p.role_id = r.id ORDER BY p.permission_code
)::text[] AS codes
FROM roles r
WHERE r.id = $1;

-- GetRoleForUpdate
SELECT r.*, ARRAY(
    SELECT p.permission_code FROM role_permissions p
    WHERE p.role_id = r.id ORDER BY p.permission_code
)::text[] AS codes
FROM roles r
WHERE r.id = $1
FOR UPDATE;

-- InsertRole
INSERT INTO roles (name, scope_type, scope_id, version, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) /* Rows */ RETURNING *;

-- UpdateRole
UPDATE roles
SET name = $1, version = $2, updated_at = $3
WHERE id = $4 AND version = $5 /* ExpectedVersion */;

-- DeleteRole
DELETE FROM roles WHERE id = $1;

-- InsertRolePermission
INSERT INTO role_permissions (role_id, permission_code) VALUES ($1, $2) /* Rows */;

-- ClearRolePermissions
DELETE FROM role_permissions WHERE role_id = $1;

-- InsertOutboxMessage
INSERT INTO outbox_messages (event_type, aggregate_id, payload) VALUES ($1, $2, $3) /* Rows */;

-- CountRoles
SELECT count(*) FROM roles r WHERE r.scope_type=$1 AND r.scope_id=$2 AND strpos(lower(r.name),lower($3::text /* Search */))>0 AND ($4::text /* PermissionCode */='' OR EXISTS(SELECT 1 FROM role_permissions rp WHERE rp.role_id=r.id AND rp.permission_code=$4));
-- PageRoles
SELECT r.*,ARRAY(SELECT p.permission_code FROM role_permissions p WHERE p.role_id=r.id ORDER BY p.permission_code)::text[] AS codes FROM roles r WHERE r.scope_type=$1 AND r.scope_id=$2 AND strpos(lower(r.name),lower($3::text /* Search */))>0 AND ($4::text /* PermissionCode */='' OR EXISTS(SELECT 1 FROM role_permissions rp WHERE rp.role_id=r.id AND rp.permission_code=$4))
ORDER BY CASE WHEN $5::text /* Sort */='name' AND $6::text /* Direction */='asc' THEN r.name END asc,
CASE WHEN $5::text /* Sort */='name' AND $6::text /* Direction */='desc' THEN r.name END desc,
CASE WHEN $5::text /* Sort */='created_at' AND $6::text /* Direction */='asc' THEN r.created_at END asc,
CASE WHEN $5::text /* Sort */='created_at' AND $6::text /* Direction */='desc' THEN r.created_at END desc,
CASE WHEN $5::text /* Sort */='updated_at' AND $6::text /* Direction */='asc' THEN r.updated_at END asc,
CASE WHEN $5::text /* Sort */='updated_at' AND $6::text /* Direction */='desc' THEN r.updated_at END desc,
CASE WHEN $6='asc' THEN r.id END ASC,
CASE WHEN $6='desc' THEN r.id END DESC
LIMIT $7::int /* Limit */ OFFSET $8::int /* Offset */;

