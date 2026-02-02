package models

import (
	"time"
)

type Task struct {
	ID          string     `json:"id" validate:"required"`
	Title       string     `json:"title" validate:"required, min=3, max=100"`
	Description *string    `json:"description" validate:"max=500"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date" validate:"required"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskUpdate struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
	Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskFilters struct {
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=pending in_progress completed"`
	Priority *string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	Page     int     `json:"page" validate:"omitempty,min=1"`
	Limit    int     `json:"limit" validate:"omitempty,min=1,max=100"`
}
