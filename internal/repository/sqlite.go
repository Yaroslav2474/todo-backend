package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	".phonepc_link/projecrt/to-do_list/todo-backend/internal/models"
)

func nullableString(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}

type sqliteRepo struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *sqliteRepo {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) Create(ctx context.Context, task *models.Task) error {
	const query = `INSERT INTO tasks (title, description, completed, priority, due_date) VALUES (?, ?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, task.Title,
		nullableString(task.Description), task.Status,
		task.Priority, nullableTime(task.DueDate))

	if err != nil {
		return fmt.Errorf("ошибка в создании задачи: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("не удалось получить последний идентификатор вставки: %w", err)
	}

	task.ID = strconv.FormatInt(id, 10)

	var createdAt, updatedAt time.Time
	err = r.db.QueryRowContext(ctx, "SELECT created_at, updated_at FROM tasks WHERE id = ?", task.ID).Scan(&createdAt, &updatedAt)
	if err != nil {
		return fmt.Errorf("не удалось получить метки времени: %w", err)
	}

	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	return nil
}

func (r *sqliteRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	const query = `
        SELECT id, title, description, status, priority, due_date, created_at, updated_at
        FROM tasks
        WHERE id = ?
    `
	task := &models.Task{}
	var duesc sql.NullString
	var DueDate sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&duesc,
		&task.Status,
		&task.Priority,
		*&DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найден: %w", err)
		}
		return nil, fmt.Errorf("не удалось получить задание: %w", err)
	}

	if duesc.Valid {
		task.Description = &duesc.String
	}
	if DueDate.Valid {
		parsedTime, err := time.Parse(time.RFC3339, DueDate.String)
		if err != nil {
			return nil, fmt.Errorf("не удалось распарсить дату: %w", err)
		}
		task.DueDate = &parsedTime
	}

	return task, nil

}

func (r *sqliteRepo) GetAll(ctx context.Context, filters map[string]string) ([]models.Task, error) {
	query := `SEKECT id, title, description, status, priority, due_date, created_at, updated_at FROM tasks`

	args := []interface{}{}

	if status, ok := filters["status"]; ok {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priority, ok := filters["priority"]; ok {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	query += "ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос к задачам: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		var duesc sql.NullString
		var dueDate sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&duesc,
			&task.Status,
			&task.Priority,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка в записи файла: %w", err)
		}

		if duesc.Valid {
			task.Description = &duesc.String
		}
		if dueDate.Valid {
			pasrsedTime, err := time.Parse(time.RFC3339, dueDate.String)
			if err == nil {
				task.DueDate = &pasrsedTime
			}
		}

	}

}
