package repository

import (
	"context"

	".phonepc_link/projecrt/to-do_list/todo-backend/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error

	GetByID(ctx context.Context, id int) (*models.Task, error)

	GetAll(ctx context.Context, filters map[string]string) ([]models.Task, error)

	Update(ctx context.Context, id int, task *models.Task) error

	Delete(ctx context.Context, id int) error
}
