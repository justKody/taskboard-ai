-- name: CreateOrganizationInvite :one
INSERT INTO organization_invites (organization_id, user_id, invited_by)
VALUES ($1, $2, $3)
RETURNING id, organization_id, user_id, invited_by, status, created_at;

-- name: GetPendingOrganizationInvite :one
SELECT id, organization_id, user_id, invited_by, status, created_at
FROM organization_invites
WHERE organization_id = $1 AND user_id = $2 AND status = 'pending';

-- name: GetOrganizationInviteById :one
SELECT id, organization_id, user_id, invited_by, status, created_at
FROM organization_invites
WHERE id = $1;

-- name: UpdateOrganizationInviteStatus :one
UPDATE organization_invites
SET status = $2
WHERE id = $1 AND status = 'pending'
RETURNING id, organization_id, user_id, invited_by, status, created_at;
