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

-- GetTemplate
SELECT r.*, ARRAY(
    SELECT p.permission_code FROM template_permissions p
    WHERE p.template_id = r.id ORDER BY p.permission_code
)::text[] AS codes
FROM role_templates r
WHERE r.id = $1;

-- GetTemplateForUpdate
SELECT r.*, ARRAY(
    SELECT p.permission_code FROM template_permissions p
    WHERE p.template_id = r.id ORDER BY p.permission_code
)::text[] AS codes
FROM role_templates r
WHERE r.id = $1
FOR UPDATE;

-- InsertTemplate
INSERT INTO role_templates (name, target_scope_type, system_id, version, created_by, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7) /* Rows */ RETURNING *;

-- UpdateTemplate
UPDATE role_templates
SET name = $1, version = $2, updated_at = $3
WHERE id = $4 AND version = $5 /* ExpectedVersion */;

-- DeleteTemplate
DELETE FROM role_templates WHERE id = $1;

-- InsertTemplatePermission
INSERT INTO template_permissions (template_id, permission_code) VALUES ($1, $2) /* Rows */;

-- ClearTemplatePermissions
DELETE FROM template_permissions WHERE template_id = $1;

-- LockAssignmentRole
SELECT id FROM roles WHERE id = $1 FOR UPDATE;

-- InsertAssignment
INSERT INTO role_assignments (role_id, user_id, assigned_by, assigned_at)
VALUES ($1, $2, $3, $4) /* Rows */;

-- DeleteAssignment
DELETE FROM role_assignments WHERE role_id = $1 AND user_id = $2;

-- ListAssignments
SELECT user_id, assigned_by, assigned_at
FROM role_assignments WHERE role_id = $1 ORDER BY user_id;

-- InsertOutboxMessage
INSERT INTO outbox_messages (event_type, aggregate_id, payload) VALUES ($1, $2, $3) /* Rows */;

-- GetAssignment
SELECT role_id,user_id,assigned_by,assigned_at FROM role_assignments WHERE role_id=$1 AND user_id=$2;

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

-- CountTemplates
SELECT count(*) FROM role_templates r WHERE r.target_scope_type=$1 AND r.system_id=$2 AND strpos(lower(r.name),lower($3::text /* Search */))>0;
-- PageTemplates
SELECT r.*,ARRAY(SELECT p.permission_code FROM template_permissions p WHERE p.template_id=r.id ORDER BY p.permission_code)::text[] AS codes FROM role_templates r WHERE r.target_scope_type=$1 AND r.system_id=$2 AND strpos(lower(r.name),lower($3::text /* Search */))>0
ORDER BY CASE WHEN $4::text /* Sort */='name' AND $5::text /* Direction */='asc' THEN r.name END asc,
CASE WHEN $4::text /* Sort */='name' AND $5::text /* Direction */='desc' THEN r.name END desc,
CASE WHEN $4::text /* Sort */='created_at' AND $5::text /* Direction */='asc' THEN r.created_at END asc,
CASE WHEN $4::text /* Sort */='created_at' AND $5::text /* Direction */='desc' THEN r.created_at END desc,
CASE WHEN $4::text /* Sort */='updated_at' AND $5::text /* Direction */='asc' THEN r.updated_at END asc,
CASE WHEN $4::text /* Sort */='updated_at' AND $5::text /* Direction */='desc' THEN r.updated_at END desc,
CASE WHEN $5='asc' THEN r.id END ASC,
CASE WHEN $5='desc' THEN r.id END DESC
LIMIT $6::int /* Limit */ OFFSET $7::int /* Offset */;

-- CountAssignments
SELECT count(*) FROM role_assignments a JOIN roles r ON r.id=a.role_id WHERE ($1::int /* RoleID */=0 OR a.role_id=$1) AND ($2::int /* UserID */=0 OR a.user_id=$2) AND ($3::text /* ScopeType */='' OR r.scope_type=$3) AND ($4::int /* ScopeID */=0 OR r.scope_id=$4);
-- PageAssignments
SELECT a.role_id,a.user_id,a.assigned_by,a.assigned_at,r.name AS role_name,r.scope_type,r.scope_id FROM role_assignments a JOIN roles r ON r.id=a.role_id WHERE ($1::int /* RoleID */=0 OR a.role_id=$1) AND ($2::int /* UserID */=0 OR a.user_id=$2) AND ($3::text /* ScopeType */='' OR r.scope_type=$3) AND ($4::int /* ScopeID */=0 OR r.scope_id=$4) ORDER BY a.role_id,a.user_id LIMIT $5::int /* Limit */ OFFSET $6::int /* Offset */;
