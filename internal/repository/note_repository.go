package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-restapi/internal/model"
)

type NoteRepository interface {
	Create(ctx context.Context, note *model.Note) error
	GetById(ctx context.Context, id int) (*model.Note, error)
	GetAll(ctx context.Context) ([]model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) NoteRepository {
	return &noteRepository{db: db}
}

func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
	query := `
		INSERT INTO note (title, content)
		VALUES ($1, $2)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query, note.Title, note.Content).Scan(&note.ID)
}

func (r *noteRepository) GetById(ctx context.Context, id int) (*model.Note, error) {
	query := `
		SELECT * 
		FROM note 
		WHERE id = $1
	`
	var note model.Note
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Title, &note.Content); err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepository) GetAll(ctx context.Context) ([]model.Note, error) {
	query := `
		SELECT *
		FROM note
		ORDER BY id
		`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notes []model.Note
	for rows.Next() {
		note := model.Note{}
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil

}

// Update обновляет существующую заметку
func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
	query := `
		UPDATE note
		SET title = $1, content = $2
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, note.Title, note.Content, note.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

// Delete удаляет заметку по ID
func (r *noteRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM note
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("note not found")
	}

	return nil
}

// Count возвращает общее количество пользователей
func (r *noteRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM note`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
