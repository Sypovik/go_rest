package handler

import (
	"context"
	"database/sql"
	"go-restapi/internal/model"
)

type NoteHandler interface {
	// Create(ctx context.Context, note *model.Note) error
	GetById(ctx context.Context, id int) (*model.Note, error)
	GetAll(ctx context.Context) ([]model.Note, error)
	// Update(ctx context.Context, note *model.Note) error
	// Delete(ctx context.Context, id int) error
}

type noteRepository struct {
	db *sql.DB
}

func NewNoteHandler(db *sql.DB) NoteHandler {
	return &noteRepository{db: db}
}

func (r *noteRepository) GetById(ctx context.Context, id int) (*model.Note, error) {
	query := `
		SELECT * 
		FROM note 
		WHERE id = $1
	`
	var note model.Note
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&note.Id, &note.Title, &note.Content); err != nil {
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
		if err := rows.Scan(&note.Id, &note.Title, &note.Content); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil

}

// func NewNoteHandler(db *sql.DB) NoteHandler {
// 	return &noteRepository{db: db}
// }

// // Create создает новую заметку в БД
// func (r *noteRepository) Create(ctx context.Context, note *model.Note) error {
// 	query := `
// 		INSERT INTO note (title, content)
// 		VALUES ($1, $2)
// 		RETURNING id
// 	`

// 	err := r.db.QueryRowContext(ctx, query, note.Title, note.Content).Scan(&note.Id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // GetById возвращает заметку по Id
// func (r *noteRepository) GetById(ctx context.Context, id int) (*model.Note, error) {
// 	query := `
// 		SELECT id, title, content
// 		FROM note
// 		WHERE id = $1
// 	`

// 	var note model.Note
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(&note.Id, &note.Title, &note.Content)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errors.New("note not found")
// 		}
// 		return nil, err
// 	}

// 	return &note, nil
// }

// // GetAll возвращает все заметки
// func (r *noteRepository) GetAll(ctx context.Context) ([]model.Note, error) {
// 	query := `
// 		SELECT id, title, content
// 		FROM note
// 		ORDER BY id DESC
// 	`

// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var notes []model.Note
// 	for rows.Next() {
// 		var note model.Note
// 		if err := rows.Scan(&note.Id, &note.Title, &note.Content); err != nil {
// 			return nil, err
// 		}
// 		notes = append(notes, note)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return notes, nil
// }

// // Update обновляет существующую заметку
// func (r *noteRepository) Update(ctx context.Context, note *model.Note) error {
// 	query := `
// 		UPDATE note
// 		SET title = $1, content = $2
// 		WHERE id = $3
// 	`

// 	result, err := r.db.ExecContext(ctx, query, note.Title, note.Content, note.Id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("note not found")
// 	}

// 	return nil
// }

// // Delete удаляет заметку по Id
// func (r *noteRepository) Delete(ctx context.Context, id int) error {
// 	query := `
// 		DELETE FROM note
// 		WHERE id = $1
// 	`

// 	result, err := r.db.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return errors.New("note not found")
// 	}

// 	return nil
// }
