package service

import (
	"context"
	"database/sql"
	"errors"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"log"
)

// Определение ошибок сервисного уровня
var (
	ErrNoteNotFound   = errors.New("note not found")
	ErrInvalidNoteID  = errors.New("invalid note ID")
	ErrEmptyNoteTitle = errors.New("note title cannot be empty")
)

// NoteService определяет интерфейс сервисного слоя для заметок
type NoteService interface {
	CreateNote(ctx context.Context, note *model.Note) error
	GetNoteByID(ctx context.Context, id int) (*model.Note, error)
	GetAllNotes(ctx context.Context) ([]model.Note, error)
	UpdateNote(ctx context.Context, note *model.Note) error
	DeleteNote(ctx context.Context, id int) error
	GetNotesCount(ctx context.Context) (int, error)
}

// noteService реализует интерфейс NoteService
type noteService struct {
	repo repository.NoteRepository
}

// NewNoteService создает новый экземпляр сервиса заметок
func NewNoteService(repo repository.NoteRepository) NoteService {
	return &noteService{
		repo: repo,
	}
}

// CreateNote создает новую заметку
func (s *noteService) CreateNote(ctx context.Context, note *model.Note) error {
	// Валидация
	if note.Title == "" {
		return ErrEmptyNoteTitle
	}

	// Логирование
	log.Printf("Creating new note: %s", note.Title)

	// Вызов репозитория
	if err := s.repo.Create(ctx, note); err != nil {
		log.Printf("Error creating note: %v", err)
		return err
	}

	return nil
}

// GetNoteByID получает заметку по ID
func (s *noteService) GetNoteByID(ctx context.Context, id int) (*model.Note, error) {
	// Валидация
	if id <= 0 {
		return nil, ErrInvalidNoteID
	}

	// Вызов репозитория
	note, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		log.Printf("Error getting note by ID %d: %v", id, err)
		return nil, err
	}

	return note, nil
}

// GetAllNotes получает все заметки
func (s *noteService) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	notes, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error getting all notes: %v", err)
		return nil, err
	}

	return notes, nil
}

// UpdateNote обновляет существующую заметку
func (s *noteService) UpdateNote(ctx context.Context, note *model.Note) error {
	// Валидация
	if note.ID <= 0 {
		return ErrInvalidNoteID
	}
	if note.Title == "" {
		return ErrEmptyNoteTitle
	}

	// Проверка существования заметки
	_, err := s.repo.GetById(ctx, note.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoteNotFound
		}
		return err
	}

	// Логирование
	log.Printf("Updating note ID %d: %s", note.ID, note.Title)

	// Вызов репозитория
	if err := s.repo.Update(ctx, note); err != nil {
		log.Printf("Error updating note ID %d: %v", note.ID, err)
		return err
	}

	return nil
}

// DeleteNote удаляет заметку по ID
func (s *noteService) DeleteNote(ctx context.Context, id int) error {
	// Валидация
	if id <= 0 {
		return ErrInvalidNoteID
	}

	// Проверка существования заметки
	_, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoteNotFound
		}
		return err
	}

	// Логирование
	log.Printf("Deleting note ID %d", id)

	// Вызов репозитория
	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting note ID %d: %v", id, err)
		return err
	}

	return nil
}

// GetNotesCount возвращает количество заметок
func (s *noteService) GetNotesCount(ctx context.Context) (int, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		log.Printf("Error counting notes: %v", err)
		return 0, err
	}

	return count, nil
}
