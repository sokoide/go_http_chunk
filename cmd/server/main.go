package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		reader := bufio.NewReader(r.Body)
		buffer := make([]byte, 1024)

		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				logrus.Infof("received: %s\n", string(buffer[:n]))
			}

			if err == io.EOF {
				break
			}
			if err != nil {
				logrus.Errorf("err: %v", err)
				break
			}
		}
		logrus.Info("Done")
		fmt.Fprintf(w, "Data Received\n")
	} else {
		http.ServeFile(w, r, "form.html")
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000",
	})

	server := &http.Server{
		Addr:           ":8080",
		MaxHeaderBytes: 1 << 20, // 1MB
		Handler:        http.HandlerFunc(handler),
	}

	logrus.Info("server starting...")
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatal(err)
	}
}
