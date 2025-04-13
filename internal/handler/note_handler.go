package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"go-restapi/internal/model"
	"go-restapi/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteHandler interface {
	RegisterRoutes(router *gin.Engine)
	GetNoteByID(c *gin.Context)
	GetAllNotes(c *gin.Context)
	CreateNote(c *gin.Context)
	UpdateNote(c *gin.Context)
	DeleteNote(c *gin.Context)
	GetNotesCount(c *gin.Context)
}

type noteHandler struct {
	service service.NoteService
}

func NewNoteHandler(service service.NoteService) NoteHandler {
	return &noteHandler{
		service: service,
	}
}

// RegisterRoutes регистрирует маршруты для заметок в Gin engine
func (h *noteHandler) RegisterRoutes(router *gin.Engine) {
	notes := router.Group("/")
	{
		notes.POST("", h.CreateNote)
		notes.GET("", h.GetAllNotes)
		notes.GET("/:id", h.GetNoteByID)
		notes.PUT("/:id", h.UpdateNote)
		notes.DELETE("/:id", h.DeleteNote)
		notes.GET("/count", h.GetNotesCount)
	}
}

func (h *noteHandler) CreateNote(c *gin.Context) {
	var note model.Note

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.service.CreateNote(c.Request.Context(), &note)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmptyNoteTitle):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *noteHandler) UpdateNote(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	note.Id = id

	if err := h.service.UpdateNote(c.Request.Context(), &note); err != nil {
		switch err {
		case service.ErrNoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		case service.ErrInvalidNoteID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		case service.ErrEmptyNoteTitle:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Note title cannot be empty"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		}
		return
	}

	c.JSON(http.StatusCreated, note)
}

func (h *noteHandler) DeleteNote(c *gin.Context) {
	// 1. Извлечение и валидация параметра ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	if err := h.service.DeleteNote(c.Request.Context(), id); err != nil {
		switch err {
		case service.ErrNoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		case service.ErrInvalidNoteID:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *noteHandler) GetAllNotes(c *gin.Context) {
	notes, err := h.service.GetAllNotes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *noteHandler) GetNoteByID(c *gin.Context) {
	// 1. Извлечение и валидация параметра ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		return
	}

	// 2. Получение заметки через сервисный слой
	note, err := h.service.GetNoteByID(c.Request.Context(), id)
	log.Println("Функция GetNoteByID запущена внутри handler")
	log.Printf("ID заметки: %d", id)
	if err != nil {
		// 3. Обработка различных типов ошибок
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		case errors.Is(err, service.ErrInvalidNoteID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID format"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	log.Printf("Заметка получена внутри handler: %d", note.Id)

	// 4. Успешный ответ
	c.JSON(http.StatusOK, note)
}

func (h *noteHandler) GetNotesCount(c *gin.Context) {
	count, err := h.service.GetNotesCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка в получении count"})
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
