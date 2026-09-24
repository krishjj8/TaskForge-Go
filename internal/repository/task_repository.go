package repository

import (
	"context"
	"database/sql"
	"errors"

	"taskforge/internal/model"
)

// ErrNotFound is returned when a requested task does not exist in storage
var ErrNotFound = errors.New("task not found")

// TaskRepository defines all storage operations for tasks
type TaskRepository interface {
	GetAll(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id int) (*model.Task, error)
	Create(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, id int, task *model.Task) error
	Delete(ctx context.Context, id int) error
}

// PostgresRepository implements TaskRepository using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository returns a new instance of PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	query := `SELECT id, title, COALESCE(description, ''), done FROM tasks ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*model.Task, error) {
	query := `SELECT id, title, COALESCE(description, ''), done FROM tasks WHERE id = $1`

	var t model.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Done)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (r *PostgresRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (title, description, done)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query, task.Title, task.Description, task.Done).Scan(&task.ID)
}

func (r *PostgresRepository) Update(ctx context.Context, id int, task *model.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, done = $3
		WHERE id = $4
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query, task.Title, task.Description, task.Done, id).Scan(&task.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
