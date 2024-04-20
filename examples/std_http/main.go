package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		w.Write([]byte("OK 1"))
	})

	s := http.Server{
		Addr:    ":7070",
		Handler: mux,
	}

	log.Println(s.ListenAndServe())

	s.Shutdown(nil)
}
