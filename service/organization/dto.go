package organization

type CreateOrganizationRequestDTO struct {
	Name string `json:"name" validate:"required"`
}
