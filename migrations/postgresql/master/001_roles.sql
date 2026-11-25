-- +goose Up
CREATE TABLE roles (
 id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name TEXT NOT NULL CHECK (btrim(name) <> ''),
 scope_type TEXT NOT NULL CHECK (scope_type IN ('system','workspace','project')),
 scope_id INT NOT NULL, version BIGINT NOT NULL CHECK (version > 0),
 created_by INT NOT NULL, created_at TIMESTAMPTZ NOT NULL, updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX roles_scope ON roles(scope_type,scope_id,id);
CREATE TABLE role_permissions (
 role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
 permission_code TEXT NOT NULL, PRIMARY KEY(role_id,permission_code)
);
CREATE TABLE role_assignments (
 role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
 user_id INT NOT NULL, assigned_by INT NOT NULL, assigned_at TIMESTAMPTZ NOT NULL,
 PRIMARY KEY(role_id,user_id)
);
CREATE TABLE role_templates (
 id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, system_id INT NOT NULL, name TEXT NOT NULL CHECK (btrim(name) <> ''),
 target_scope_type TEXT NOT NULL CHECK (target_scope_type IN ('system','workspace','project')),
 version BIGINT NOT NULL CHECK (version > 0),
 created_by INT NOT NULL, created_at TIMESTAMPTZ NOT NULL, updated_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX templates_scope ON role_templates(system_id,target_scope_type,id);
CREATE TABLE template_permissions (
 template_id INT NOT NULL REFERENCES role_templates(id) ON DELETE CASCADE,
 permission_code TEXT NOT NULL, PRIMARY KEY(template_id,permission_code)
);
CREATE TABLE outbox_messages (
 id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 event_type TEXT NOT NULL, aggregate_id INT NOT NULL,
 payload JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE outbox_messages;
DROP TABLE template_permissions;
DROP TABLE role_templates;
DROP TABLE role_assignments;
DROP TABLE role_permissions;
DROP TABLE roles;
