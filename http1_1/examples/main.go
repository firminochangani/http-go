package main

import (
	"encoding/json"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"http_v1_1"
)

func main() {
	router := http1_1.NewServerDefaultRouter()
	router.GET("/", func(r *http1_1.Request, w *http1_1.Response) error {
		log.Println("handling request")
		w.Headers.Set("Content-Type", "text/html; charset=UTF-8")
		return w.Write([]byte("Hello world"))
	})

	router.GET("/peoples", func(r *http1_1.Request, w *http1_1.Response) error {
		peoples := make([]map[string]string, 100)

		for i := 0; i < 100; i++ {
			peoples[i] = map[string]string{
				"id":         gofakeit.UUID(),
				"first_name": gofakeit.FirstName(),
				"last_name":  gofakeit.LastName(),
				"email":      gofakeit.Email(),
			}
		}

		payload, err := json.Marshal(peoples)
		if err != nil {
			w.WriteStatus(500)
			return w.Write([]byte("Internal Server Error"))
		}

		w.Headers.Set("Content-Type", "application/json; charset=UTF-8")
		return w.Write(payload)
	})

	s := http1_1.Server{
		Addr:   ":8080",
		Router: router,
	}

	log.Println(s.ListenAndServe())
}
