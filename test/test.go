package main

import (
	"context"
	"fmt"
	"go-restapi/internal/config"
	"time"
)

func printAfterDelay(ctx context.Context, delay time.Duration, text string) {
	select {
	case <-time.After(delay): // Ждем указанное время
		fmt.Println(text)
	case <-ctx.Done(): // Срабатывает при отмене контекста
		fmt.Println("Операция отменена:", ctx.Err())
	}
}

func main() {
	config := config.LoadConfig()

	Timeout := config.DBTimeout

	timeoutDuration, err := time.ParseDuration(Timeout)
	if err != nil {
		panic(fmt.Sprintf("Invalid timeout format: %v", err))
	}

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel() // Важно: освобождаем ресурсы

	fmt.Println("Ожидаем 5 секунд...")
	printAfterDelay(ctx, 5*time.Second, "Текст через 5 секунд!")

	// Для демонстрации отмены раскомментируйте:
	// go func() {
	//   time.Sleep(2*time.Second)
	//   cancel() // Принудительная отмена
	// }()
}
