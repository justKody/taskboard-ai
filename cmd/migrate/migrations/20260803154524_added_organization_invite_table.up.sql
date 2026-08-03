CREATE TYPE organization_invite_status AS ENUM (
    'pending',
    'accepted',
    'rejected'
);

CREATE TABLE organization_invites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    user_id UUID NOT NULL REFERENCES users(id),
    invited_by UUID NOT NULL REFERENCES users(id),
    status organization_invite_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- one open invite per user per organization
CREATE UNIQUE INDEX organization_invites_pending_unique
    ON organization_invites (organization_id, user_id)
    WHERE status = 'pending';
