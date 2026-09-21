package project

type CreateProjectRequestDTO struct {
	OrganizationID string `json:"organization_id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description"`
}
