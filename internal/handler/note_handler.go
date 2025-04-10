package handler

import (
	"go-restapi/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetNoteByID(c *gin.Context) {
	// Конвертация параметра ID из строки в целое число
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Возвращение ошибки 400 (Bad Request), если ID некорректен
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}
	// Получение заметки по ID из хранилища
	var note *model.Note = &model.Note{
		Id:      id,
		Title:   "add",
		Content: "asd",
	}

	if note == nil {
		// Возвращение ошибки 404 (Not Found), если заметка не найдена
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	// Возвращение найденной заметки в формате JSON с кодом 200 (OK)
	c.JSON(http.StatusOK, note)

}
