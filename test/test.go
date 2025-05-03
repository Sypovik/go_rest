package main

import (
	"context"
	"errors"
	"fmt"
	"go-restapi/internal/config"
	"math/rand"
	"sync"
	"time"
)

func testTimeout() {

	printAfterDelay := func(ctx context.Context, delay time.Duration, text string) {
		select {
		case <-time.After(delay): // Ждем указанное время
			fmt.Println(text)
		case <-ctx.Done(): // Срабатывает при отмене контекста
			fmt.Println("Операция отменена:", ctx.Err())
		}
	}

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

func testGorutings() {
	imtationWork := func(intCh chan int, delay time.Duration) {
		fmt.Println("start function c паузой в ", int(delay.Seconds()), "s")
		now := time.Now()
		time.Sleep(delay)
		intCh <- int(delay.Seconds())
		fmt.Println("stop function c паузой в ", int(delay.Seconds()), "s")
		fmt.Println("Прошло: ", int(time.Since(now).Seconds()), " секунд\n ")
	}

	intCh := make(chan int)
	for i := range 3 {
		delay, _ := time.ParseDuration(fmt.Sprintf("%ds", i+1))
		go imtationWork(intCh, delay)
	}

	fmt.Println("Запущены асинхронные функции, ожидаем результаты")
	time.Sleep(5 * time.Second)
	for range 3 {
		fmt.Println("Результат: ", <-intCh)
	}
	time.Sleep(1 * time.Second)
}

func testSync() {

	counter := 0
	work := func(number int, ch chan bool) {
		counter = 0
		for k := 1; k <= 5; k++ {
			counter++
			fmt.Println("Goroutine", number, "-", counter)
		}
		ch <- true
	}

	ch := make(chan bool) // канал
	for i := 1; i < 5; i++ {
		go work(i, ch)
	}
	// ожидаем завершения всех горутин
	for i := 1; i < 5; i++ {
		<-ch
	}
	fmt.Println("The End")
}

func task1() {

	// timeoutLimit - вероятность, с которой не будет возвращаться ошибка от fakeDownload():
	// timeoutLimit = 100 - ошибок не будет;
	// timeoutLImit = 0 - всегда будет возвращаться ошибка.
	// Можете "поиграть" с этим параметром, для проверки случаев с возвращением ошибки.
	const timeoutLimit = 90
	type Result struct {
		msg string
		err error
	}
	// fakeDownload - имитирует разное время скачивания для разных адресов
	fakeDownload := func(url string) Result {
		r := rand.Intn(100)
		time.Sleep(time.Duration(r) * time.Millisecond)
		if r > timeoutLimit {
			return Result{
				err: fmt.Errorf("failed to download data from %s: timeout", url),
			}
		}
		return Result{
			msg: fmt.Sprintf("downloaded data from %s\n", url),
		}
	}
	// download - параллельно скачивает данные из urls
	download := func(urls []string) ([]string, error) {

		var errs error
		var msgs []string
		var wg sync.WaitGroup
		wg.Add(len(urls))
		c := make(chan Result, len(urls))
		for _, url := range urls {
			go func(url string) {
				defer wg.Done()
				c <- fakeDownload(url)
			}(url)
		}

		wg.Wait()
		close(c)

		for r := range c {
			if r.err != nil {
				errs = errors.Join(r.err)
			}
			msgs = append(msgs, r.msg)
		}

		return msgs, errs
	}

	msgs, err := download([]string{
		"https://example.com/e25e26d3-6aa3-4d79-9ab4-fc9b71103a8c.xml",
		"https://example.com/a601590e-31c1-424a-8ccc-decf5b35c0f6.xml",
		"https://example.com/1cf0dd69-a3e5-4682-84e3-dfe22ca771f4.xml",
		"https://example.com/ceb566f2-a234-4cb8-9466-4a26f1363aa8.xml",
		"https://example.com/b6ed16d7-cb3d-4cba-b81a-01a789d3a914.xml",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(msgs)
}

func task2() {
	repeatFn := func(ctx context.Context, fn func() any) <-chan any {
		out := make(chan any)
		go func() {
			defer close(out)
			for {
				select {
				case <-ctx.Done():
					return
				case out <- fn():
				}
			}
		}()
		return out
	}

	take := func(ctx context.Context, in <-chan any, num int) <-chan any {
		out := make(chan any)
		go func() {
			defer close(out)
			for range num {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					out <- v
				}
			}
		}()
		return out
	}
	main := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		rand := func() interface{} { return rand.Int() }
		var res []interface{}
		for num := range take(ctx, repeatFn(ctx, rand), 3) {
			res = append(res, num)
		}
		if len(res) != 3 {
			panic("wrong code")
		}

	}
	main()
}

func task3() {
	worker := func(f func(int) int, jobs <-chan int, results chan<- int) {
		for job := range jobs {
			results <- f(job)
		}
	}
	const numJobs = 5
	const numWorkers = 3
	main := func() {
		jobs := make(chan int, numJobs)
		results := make(chan int, numJobs)
		wg := sync.WaitGroup{}
		multiplier := func(x int) int {
			return x * 10
		}
		wg.Add(numWorkers)
		for range numWorkers {
			go func() {
				defer wg.Done()
				worker(multiplier, jobs, results)
			}()
		}
		done := make(chan struct{})
		var arrInt []int
		go func() {
			for r := range results {
				arrInt = append(arrInt, r)
			}
			close(done)
		}()
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
		go func() {
			wg.Wait()
			close(results)
		}()
		<-done

	}
	main()

}

func task4() {
	mergeSorted := func(a, b <-chan int) <-chan int {
		out := make(chan int)
		go func() {

			for i := range a {
				out <- i
			}

			_, ok := <-b
			if !ok {

				close(out)
			}
		}()
		go func() {

			for i := range b {
				out <- i
			}

			_, ok := <-a
			if !ok {
				close(out)
			}
		}()
		return out
	}
	fillChanA := func(c chan int) {
		c <- 1
		c <- 2
		c <- 4
		close(c)
	}
	fillChanB := func(c chan int) {
		c <- -1
		c <- 4
		c <- 5
		close(c)
	}
	main := func() {
		a, b := make(chan int), make(chan int)
		go fillChanA(a)
		go fillChanB(b)
		c := mergeSorted(a, b)
		go func() {

			for val := range c {
				fmt.Printf("%d ", val)
			}
		}()
	}
	main()
}

func task5() {
	// merge - соединяет каналы в один
	merge := func(cs ...<-chan int) <-chan int {
		result := make(chan int)
		var wg sync.WaitGroup
		wg.Add(len(cs))
		for _, ch := range cs {
			go func() {
				for c := range ch {
					result <- c
				}
				wg.Done()
			}()
		}
		go func() {
			wg.Wait()
			close(result)
		}()
		return result
	}
	// fillChan - заполняет канал числами от 0 до n-1
	fillChan := func(n int) <-chan int {
		result := make(chan int)
		go func() {
			for i := range n {
				result <- i
			}
			close(result)
		}()
		return result
	}
	main := func() {
		a := fillChan(2)
		b := fillChan(3)
		c := fillChan(4)
		d := merge(a, b, c)
		for v := range d {
			fmt.Println(v)
		}
	}
	main()

}

func task6() {
	mergeSorted := func(a, b <-chan int) <-chan int {
		result := make(chan int)
		go func() {
			defer close(result)
			val1, ok1 := <-a
			val2, ok2 := <-b
			for ok1 && ok2 {
				if val1 < val2 {
					result <- val1
					// result <- val2
					val1, ok1 = <-a
				} else {
					result <- val2
					// result <- val1
					val2, ok2 = <-b
				}
			}
			for ok1 {
				result <- val1
				// result <- val2
				val1, ok1 = <-a
			}
			for ok2 {
				result <- val2
				// result <- val1
				val2, ok2 = <-b
			}
		}()
		return result
	}
	fillChanA := func(c chan int) {
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}
	fillChanB := func(c chan int) {
		c <- 4
		c <- 5
		c <- 6
		close(c)
	}
	main := func() {
		a, b := make(chan int), make(chan int)
		go fillChanA(a)
		go fillChanB(b)
		c := mergeSorted(a, b)
		for val := range c {
			fmt.Printf("%d ", val)
		}
	}

	main()

}

func main() {
	_ = testTimeout
	_ = testGorutings
	_ = testSync
	_ = task1
	_ = task2
	_ = task3
	_ = task4
	_ = task5

	task6()
}
