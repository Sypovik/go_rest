package main

import (
	"context"
	"database/sql"
	"go-restapi/internal/config"
	"go-restapi/internal/model"
	"go-restapi/internal/repository"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func conectToDB() (*sql.DB, error) {
	config := config.LoadConfig()

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка создания объекта БД: %v", err) // Переименуем сообщение
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
	var dr repository.NoteRepository
	dr = conectToRepository()

	// Правильно сохраняем и используем функцию отмены
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Важно: всегда вызывайте cancel в конце

	// var note model.Note
	notes, err := dr.GetAll(ctx)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Полученные заметки:")
	for _, note := range notes {
		log.Printf("\tnote.Id:%d\tnote.Title:%s\tnote.Content:%s\n", note.Id, note.Title, note.Content)
	}
	// note = model.Note{
	// 	Title:   "Заметка 1",
	// 	Content: "Содержимое заметки 1",
	// }
	dr.Create(ctx, &model.Note{
		Title:   "Заметка 1",
		Content: "Содержимое заметки 1",
	})
	dr.Update(ctx, &model.Note{
		Id:      10,
		Title:   "Обновленная заметка 1",
		Content: "Обновленное содержимое заметки 1",
	})
	// Получаем обновленный список после создания
	notes, err = dr.GetAll(ctx)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	log.Println("Полученные заметки после создания:")
	for _, note := range notes {
		log.Printf("\tnote.Id:%d\tnote.Title:%s\tnote.Content:%s\n", note.Id, note.Title, note.Content)
	}
}
