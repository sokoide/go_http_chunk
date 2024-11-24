package main

import (
	"io"
	"net/http"
	"time"
)

func makeData(l int) *[]byte {
	by := make([]byte, l)
	for i := 0; i < l; i++ {
		// 41-176
		by[i] = (byte)(41 + i%(176-40))
	}
	return &by
}

func main() {
	url := "http://localhost:8081"
	pr, pw := io.Pipe()

	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Transfer-Encoding", "chunked")

	go func() {
		defer pw.Close()

		// send 500 bytes every 500ms
		for i := 0; i < 10; i++ {
			pw.Write(*makeData(1000))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
