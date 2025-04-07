package main

import (
	"context"
	"database/sql"
	"go-restapi/internal/config"
	"go-restapi/internal/handler"
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

func conectToHandler() handler.NoteHandler {
	db, _ := conectToDB()
	return handler.NewNoteHandler(db)
}

func main() {
	var dr handler.NoteHandler
	dr = conectToHandler()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	note, err := dr.GetById(ctx, 3)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}

	log.Printf("\tnote.Id:%d\tnote.Title:%s\tnote.Content:%s\n", note.Id, note.Title, note.Content)

	notes, err := dr.GetAll(ctx)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	for _, note := range notes {
		log.Printf("\tnote.Id:%d\tnote.Title:%s\tnote.Content:%s\n", note.Id, note.Title, note.Content)
	}

	// row := db.QueryRow("SELECT id, title, content FROM note WHERE id = $1", 1)
	// if err := row.Scan(&note.Id, &note.Title, &note.Content); err != nil {
	// 	log.Fatalf("Ошибка выполнения запроса: %v", err)
	// }
	// log.Printf("ID: %d, Title: %s, Content: %s\n", note.Id, note.Title, note.Content)

	// for rows.Next() {
	// 	if err := rows.Scan(&id, &name, &age); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	// }
}
