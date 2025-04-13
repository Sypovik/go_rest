package main

import (
	"context"
	"database/sql"
	"go-restapi/internal/config"
	"go-restapi/internal/handler"
	"go-restapi/internal/repository"
	"go-restapi/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func conectToDB() (*sql.DB, error) {
	config := config.LoadConfig()

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка создания объекта БД: %v", err) // Переименуем сообщение
	}

	Timeout, err := time.ParseDuration(config.DBTimeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err) // Здесь будет реальная ошибка
	}
	log.Println("Подключение к БД успешно")
	return db, nil
}

func conectToRepository() repository.NoteRepository {
	db, _ := conectToDB()
	return repository.NewNoteRepository(db)
}

func main() {
	var dr repository.NoteRepository = conectToRepository()
	var sr service.NoteService = service.NewNoteService(dr)
	var hand handler.NoteHandler = handler.NewNoteHandler(sr)
	// Правильно сохраняем и используем функцию отмены

	r := gin.Default()
	hand.RegisterRoutes(r)

	r.Run(":8083") // Запускаем сервер на порту 8083
}
