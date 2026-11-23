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

-- CountAssignments
SELECT count(*) FROM role_assignments a JOIN roles r ON r.id=a.role_id WHERE ($1::int /* RoleID */=0 OR a.role_id=$1) AND ($2::int /* UserID */=0 OR a.user_id=$2) AND ($3::text /* ScopeType */='' OR r.scope_type=$3) AND ($4::int /* ScopeID */=0 OR r.scope_id=$4);
-- PageAssignments
SELECT a.role_id,a.user_id,a.assigned_by,a.assigned_at,r.name AS role_name,r.scope_type,r.scope_id FROM role_assignments a JOIN roles r ON r.id=a.role_id WHERE ($1::int /* RoleID */=0 OR a.role_id=$1) AND ($2::int /* UserID */=0 OR a.user_id=$2) AND ($3::text /* ScopeType */='' OR r.scope_type=$3) AND ($4::int /* ScopeID */=0 OR r.scope_id=$4) ORDER BY a.role_id,a.user_id LIMIT $5::int /* Limit */ OFFSET $6::int /* Offset */;
