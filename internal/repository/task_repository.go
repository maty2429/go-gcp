package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"task-api/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, description, completed, created_at
		FROM tasks
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Create(ctx context.Context, title, description string) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description, completed, created_at
	`, title, description).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (model.Task, error) {
	var task model.Task
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, description, completed, created_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}
