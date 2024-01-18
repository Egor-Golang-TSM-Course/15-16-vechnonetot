package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func makeRequest(ctx context.Context, url string, client *http.Client, ch chan<- string) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ch <- fmt.Sprintf("Error request for %s: %v", url, err)
		return
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Error making request %s: %v", url, err)
		return
	}

	defer resp.Body.Close()

	ch <- fmt.Sprintf("Response from %s: %d", url, resp.StatusCode)
}

func main() {
	urls := []string{
		"http://example.com",
		"http://mail.ru",
		"http://yandex.by",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//http клиент
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
		},
	}
	//канал для  получения рез запросов
	resultChan := make(chan string, len(urls))

	//паралельный запуск
	for _, url := range urls {
		go makeRequest(ctx, url, client, resultChan)
	}

	//завершение всех запросов по истечение таймаутов
	for i := 0; i < len(urls); i++ {
		select {
		case result := <-resultChan:
			fmt.Println(result)
		case <-ctx.Done():
			fmt.Println("Timeout reached. Cancelling remaining requests.")
			return
		}
	}
}
