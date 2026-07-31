package organization

import "github.com/justKody/taskboard-go-api/types"

type CreateOrganizationRequestDTO struct {
	Name string `json:"name" validate:"required"`
}

type ChangeOwnerOfOrganizationRequestDTO struct {
	NewOwnerID string `json:"new_owner_id" validate:"required"`
}

type OrganizationDetailsResponseDTO struct {
	Organization *types.Organization `json:"organization"`
	Members      []types.Membership  `json:"members"`
}
