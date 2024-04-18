package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK 1"))
		w.Write([]byte("OK 2"))
		w.Header().Set("x-demo", time.Now().String())
		w.WriteHeader(http.StatusCreated)
	})

	log.Println(http.ListenAndServe(":7070", mux))
}
