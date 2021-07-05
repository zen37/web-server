//https://golang.org/doc/articles/wiki/#tmp_3

package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
	handler is of type  http.HandlerFunc
	w assembles the http server's response; by writing to it we send data to the http client
	r is the http client's request
*/
func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(w, r.URL.Path)
	fmt.Fprintf(w, "%s s-print %q q-print %v v-print \n", "Hello World", "Hallo Welt", "Hola Mundo")
	fmt.Fprintf(w, "r.RequestURI is: %s\n", r.RequestURI)
	fmt.Fprintf(w, "r.URL.Path is: %s\n", r.URL.Path)
	fmt.Fprintf(w, "r.URL is: %s\n", r.URL)
	fmt.Fprintf(w, "r.URL.User is: %s\n", r.URL.User)
	fmt.Fprintf(w, "r.URL.Host is: %s\n", r.URL.Host)
	fmt.Fprintf(w, "r.URL.Fragment is: %s\n", r.URL.Fragment)
}

func main() {

	//tells the http package to handle all requests to the web root ("/") with handler.
	http.HandleFunc("/", handler)

	// tells http to listen to port 8080, and the program will block until it is terminated
	// ListenAndServe only returns when an unexpected error occurs, so it always return an error
	log.Fatal(http.ListenAndServe(":8080", nil))
}
