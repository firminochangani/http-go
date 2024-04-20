package main

import (
	"log"
	"os"

	"github.com/flowck/http-go"
)

func main() {
	router := http_go.NewServerDefaultNaiveRouter()

	router.GET("/", func(r *http_go.Request, w *http_go.Response) error {
		w.Headers.Set("X-From", "Pure implementation of the HTTP 1.1 protocol")
		w.Headers.Set("Content-Type", "text/html; charset=UTF-8")
		return w.Write([]byte("Hello world"))
	})

	s := http_go.Server{
		Addr:   ":8080",
		Router: router,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
