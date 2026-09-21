package task

import "time"

type CreateTaskRequestDTO struct {
	ProjectID   string    `json:"project_id" validate:"required"`
	AssignedTo  string    `json:"assigned_to" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=todo in_progress done"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type UpdateTaskRequestDTO struct {
	AssignedTo  string    `json:"assigned_to" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=todo in_progress done"`
	Priority    string    `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

