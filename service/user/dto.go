package user

type SignupRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=3"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

type UpdateUserRequestDTO struct {
	Name string `json:"name" validate:"required,min=3"`
}

type ChangePasswordRequestDTO struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
