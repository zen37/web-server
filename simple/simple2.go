//https://www.calhoun.io/why-cant-i-pass-this-function-as-an-http-handler/

/*
Go has not one, but two ways to register an http handler function: http.Handle and http.HandleFunc.
The two are VERY similar. In fact, the only real difference is what you pass in as the second argument,
and even that can feel *very* similar at times.

*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

type handler struct {
	name string
}

func handlerfunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is a function with arguments: w http.ResponseWriter, r *http.Request")
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is method ServeHTTP with arguments (w http.ResponseWriter, r *http.Request) and receiver: ", h.name)
	fmt.Fprint(w, "hello, this is: ", h.name)
}

func main() {

	h := handler{name: "here's the receiver name"}
	http.Handle("/test1", h)
	http.HandleFunc("/test2", handlerfunc)
	log.Fatalln(http.ListenAndServe(":8080", nil))

}
