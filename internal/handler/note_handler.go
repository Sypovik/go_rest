package handler

import (
	"errors"
	"net/http"
	"strconv"

	"go-restapi/internal/service"

	"github.com/gin-gonic/gin"
)

type NoteHandler interface {
	GetNoteByID(c *gin.Context)
}

type noteHandler struct {
	service service.NoteService
}

func NewNoteHandler(service service.NoteService) NoteHandler {
	return &noteHandler{
		service: service,
	}
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
	if err != nil {
		// 3. Обработка различных типов ошибок
		switch {
		case errors.Is(err, service.ErrNoteNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// 4. Успешный ответ
	c.JSON(http.StatusOK, note)
}
