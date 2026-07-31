-- name: CreateOrganization :one
INSERT INTO organizations (name, owner_id)
VALUES ($1, $2)
RETURNING id, name, owner_id, created_at;

-- name: GetOrganizationById :one
SELECT id, name, owner_id, created_at FROM organizations
WHERE id = $1;

-- name: CheckIfUserIsOwner :one
SELECT owner_id FROM organizations
WHERE id = $1;


-- name: DeleteOrganizationById :exec
DELETE FROM organizations
WHERE id = $1;

-- name: ChangeOwnerOfOrganizationById :exec
UPDATE organizations
SET owner_id = $1
WHERE id = $2;

-- name: GetOrganizationsListByUserId :many
SELECT o.id, o.name, o.owner_id, o.created_at
FROM organizations o
INNER JOIN memberships m ON m.organization_id = o.id
WHERE m.user_id = $1;

