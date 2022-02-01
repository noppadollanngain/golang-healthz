package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ResponseChanel struct {
	Status    int
	Url       string
	Time      float64
	StartTime string
}

func HttpRequest(url string, ch chan ResponseChanel) {
	start := time.Now()
	resp, err := http.Get(url)
	secs := time.Since(start).Seconds()
	result := ResponseChanel{
		Url:       url,
		Time:      secs,
		StartTime: start.Format("2006-01-02 15:04:05"),
	}
	if err != nil {
		result.Status = 500
		ch <- result
		return
	}
	result.Status = resp.StatusCode
	ch <- result
}

func HandleRequest(ch chan ResponseChanel, wg *sync.WaitGroup) {
	defer wg.Done()
	result := <-ch
	reportMessage := fmt.Sprintf("%s | URL %s => status: %d, time: %f", result.StartTime, result.Url, result.Status, result.Time)
	fmt.Println(reportMessage)
}

func main() {
	urls := []string{
		"http://www.facebook.com/",
		"https://line.me/th/",
		"https://line.me/en/",
		"http://www.golang.org/",
		"https://line.me/jj/",
		"https://www.cloudflare.com/",
		"http://www.google.com/",
		"http://localhost:8000",
		"https://medium.com/",
	}

	chanel := make(chan ResponseChanel)
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go HttpRequest(url, chanel)
		go HandleRequest(chanel, &wg)
	}
	wg.Wait()
}
