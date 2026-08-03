-- name: CreateMembership :one
INSERT INTO memberships (user_id, organization_id, role)
VALUES ($1, $2, $3)
RETURNING user_id, organization_id, role, joined_at;

-- name: GetAllMembershipsByOrganizationId :many
SELECT user_id, organization_id, role, joined_at FROM memberships
WHERE organization_id = $1;

-- name: UpdateMembershipRole :exec
UPDATE memberships
SET role = $3
WHERE user_id = $1 AND organization_id = $2;