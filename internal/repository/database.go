package repository

import (
	"context"
	"database/sql"
	"go-restapi/internal/config"
	"log"
	"time"
)

func SetupDatabase() (*sql.DB, error) {
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
