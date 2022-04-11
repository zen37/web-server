package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
)

const port string = ":8080"

func main() {
	r := chi.NewRouter()

	r.Get("/en", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Get("/de", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hallo welt"))
	})

	http.ListenAndServe(port, r)

}
