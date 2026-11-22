-- +goose Up
ALTER TABLE roles ADD CONSTRAINT roles_name_length CHECK (char_length(name) <= 200);
ALTER TABLE role_templates ADD CONSTRAINT templates_name_length CHECK (char_length(name) <= 200);
-- +goose Down
ALTER TABLE role_templates DROP CONSTRAINT templates_name_length;
ALTER TABLE roles DROP CONSTRAINT roles_name_length;
