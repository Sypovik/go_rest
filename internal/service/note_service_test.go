package service

import (
	"context"
	"database/sql"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Создаем мок для интерфейса репозитория
func TestCreateNote_Success(t *testing.T) {
	// Инициализация gomock контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока репозитория
	mockRepo := repository.NewMockNoteRepository(ctrl)

	// Настройка ожидаемого вызова
	note := &model.Note{
		Title:   "Test Note",
		Content: "Test Content",
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), note).
		Return(nil)

	// Создаем сервис с моком репозитория
	service := NewNoteService(mockRepo)

	// Выполняем тестируемый метод
	err := service.CreateNote(context.Background(), note)

	// Проверяем результат
	assert.NoError(t, err)
}

func TestCreateNote_EmptyTitle(t *testing.T) {
	// Инициализация gomock контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока репозитория
	mockRepo := repository.NewMockNoteRepository(ctrl)

	// Создаем заметку с пустым заголовком
	note := &model.Note{
		Title:   "",
		Content: "Test Content",
	}

	// Создаем сервис с моком репозитория
	service := NewNoteService(mockRepo)

	// Выполняем тестируемый метод
	err := service.CreateNote(context.Background(), note)

	// Проверяем результат
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyNoteTitle, err)
}

func TestGetNoteByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockNoteRepository(ctrl)

	expectedNote := &model.Note{
		Id:      1,
		Title:   "Test Note",
		Content: "Test Content",
	}

	mockRepo.EXPECT().
		GetById(gomock.Any(), 1).
		Return(expectedNote, nil)

	service := NewNoteService(mockRepo)

	note, err := service.GetNoteByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedNote, note)
}

func TestGetNoteByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockNoteRepository(ctrl)

	mockRepo.EXPECT().
		GetById(gomock.Any(), 999).
		Return(nil, sql.ErrNoRows)

	service := NewNoteService(mockRepo)

	note, err := service.GetNoteByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, note)
	assert.Equal(t, ErrNoteNotFound, err)
}
