package main

import (
	"fmt"
	"go-restapi/internal/model"
	"io"
	"net/http" // Пакет для работы с HTTP
	"strconv"  // Пакет для конвертации строк в другие типы данных

	"github.com/gin-gonic/gin" // Веб-фреймворк Gin
)

func GetNoteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	// Конвертация параметра ID из строки в целое число
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

func EchoRequest(c *gin.Context) {
	// 1. Метод и URL
	str := fmt.Sprintf("Method: %s\n", c.Request.Method)
	str += fmt.Sprintf("URL: %s\n", c.Request.URL)

	// 2. Заголовки
	fmt.Sprintln("Headers:")
	for name, values := range c.Request.Header {
		str += fmt.Sprintf("%s: %v\n", name, values)
	}

	// 3. Query-параметры (из URL)
	str += fmt.Sprintln("Query Params:")
	for key, values := range c.Request.URL.Query() {
		str += fmt.Sprintf("%s: %v\n", key, values)
	}

	// 4. Параметры пути (например, /users/:id)
	str += fmt.Sprintln("Path Params:")
	params := c.Params
	for _, param := range params {
		str += fmt.Sprintf("%s: %s\n", param.Key, param.Value)
	}

	// 5. Тело запроса (если есть)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		str += fmt.Sprintln("Error reading body:", err)
		return
	}
	defer c.Request.Body.Close() // Важно закрыть тело!

	str += fmt.Sprintln("Body:")
	str += fmt.Sprintln(string(body))
	c.Data(200, "text/plain", []byte(str))
}

func PrintRequest(c *gin.Context) {
	// 1. Метод и URL
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("URL: %s\n", c.Request.URL)

	// 2. Заголовки
	fmt.Println("Headers:")
	for name, values := range c.Request.Header {
		fmt.Printf("%s: %v\n", name, values)
	}

	// 3. Query-параметры (из URL)
	fmt.Println("Query Params:")
	for key, values := range c.Request.URL.Query() {
		fmt.Printf("%s: %v\n", key, values)
	}

	// 4. Параметры пути (например, /users/:id)
	fmt.Println("Path Params:")
	params := c.Params
	for _, param := range params {
		fmt.Printf("%s: %s\n", param.Key, param.Value)
	}

	// 5. Тело запроса (если есть)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	defer c.Request.Body.Close() // Важно закрыть тело!

	fmt.Println("Body:")
	fmt.Println(string(body))
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Выводим информацию о запросе
		PrintRequest(c)

		// Продолжаем обработку
		c.Next()
	}
}

func main() {
	// Режим продакшена
	gin.SetMode(gin.ReleaseMode)

	// Создаем движок с доверенными прокси
	r := gin.Default()
	r.Use(RequestLogger()) // Логируем все запросы
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Роуты
	// r.GET("/notes/:id", GetNoteByID)
	r.GET("/*path", EchoRequest)

	// Запуск
	r.Run(":8080")
}
