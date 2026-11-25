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

