package types

import "time"

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// one to one or many

type Organization struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" validate:"required"`
	OwnerId   string    `json:"owner_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// many to many

type Membership struct {
	OrganizationId string `json:"organization_id"`
	UserId         string `json:"user_id"`
	Role           string `json:"role"`
	JoinedAt       string `json:"joined_at"`
}

// org → user invite (pending until accepted)
type OrganizationInvite struct {
	Id             string    `json:"id"`
	OrganizationId string    `json:"organization_id" validate:"required"`
	UserId         string    `json:"user_id" validate:"required"` // invitee
	InvitedBy      string    `json:"invited_by" validate:"required"`
	Status         string    `json:"status"` // pending/accepted/rejected
	CreatedAt      time.Time `json:"created_at"`
}

// one to many

type Project struct {
	Id             string    `json:"id"`
	OrganizationID string    `json:"organization_id" validate:"required"` // fk key
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	Status         string    `json:"status" validate:"required"` // active/completed/archived
	CreatedBy      string    `json:"created_by"`                 // fk key users // can be evalutaed by token
	CreatedAt      time.Time `json:"created_at"`
}

// tasks

type Task struct {
	Id          string    `json:"id"`
	ProjectID   string    `json:"project_id" validate:"required"` // fk
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`   // todo/in_progress/done
	Priority    string    `json:"priority"` // low/medium/high
	Due_date    string    `json:"due_date"`
	AssignedTo  string    `json:"assigned_to"` // fk
	CreatedBy   string    `json:"created_by"`  // fk
	CreatedAt   time.Time `json:"created_at"`
}

type Comments struct {
	Id       string    `json:"id"`
	TaskId   string    `json:"task_id" validate:"required"`
	UserId   string    `json:"user_id" validate:"required"`
	Message  string    `json:"message" validate:"required"`
	CreateAt time.Time `json:"created_at"`
}

type Label struct {
	Id             string `json:"id"`
	OrganizationId string `json:"organization_id" validate:"required"`
	Name           string `json:"name" validate:"required"`
	Color          string `json:"color"` // by frontend
}

// many to many

// we might not even need json format eg
type TaskLabel struct {
	TaskId  string `json:"task_id"`
	LabelId string `json:"label_id"`
}
