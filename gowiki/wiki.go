package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//describes how the wiki page will be stored in memory
type page struct {
	title string
	body  []byte // byte slice rather than string as this is the type expected by the io libraries
}

/*
save method that takes as its receiver a pointer to page
it returns an error value because that is the return type of WriteFile
if there is no error it returns nil (the zero-value for pointers, interfaces, and some other types)
octal integer literal 0600, indicates that the file should be created with read-write permissions for the current user only.
*/
func (p *page) save() error {
	f := p.title + ".txt"
	return ioutil.WriteFile(f, p.body, 0600)
}

// why not a method? tutorial has a function - May 22
func loadPage(title string) (*page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{title: title, body: body}, nil
}

// load method not in the tutorial - May 22
func (p *page) load() (*page, error) {
	f := p.title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{title: p.title, body: body}, nil

}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	/*
		var p1 page
		p1.title = r.URL.Path[len("/view/"):]
		p, err := p1.load()
	*/

	t := r.URL.Path[len("/view/"):]
	p, err := loadPage(t)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Fprintln(w, string(p2.body))
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.title, p.body)

}

//The function editHandler loads the page
//(or, if it doesn't exist, create an empty Page struct), and displays an HTML form.
func editHandler1(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Path[len("/edit1/"):]
	p1 := &page{title: t}
	p2, err := p1.load()

	if err != nil {
		fmt.Fprintf(w, "<h1>editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"body\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form>",
			p1.title, p1.title, p1.body)
	} else {
		fmt.Fprintf(w, "<h1>editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"body\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form>",
			p2.title, p2.title, p2.body)
	}
	fmt.Fprintf(w, "<h6>func (p *page) load() (*page, error) has been used</h6>")
}

func editHandler2(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Path[len("/edit2/"):]
	p, err := loadPage(t)
	if err != nil {
		p = &page{title: t}
	}

	fmt.Fprintf(w, "<h1>editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.title, p.title, p.body)
	fmt.Fprintf(w, "<h6>func loadPage(title string) (*page, error) has been used</h6>")
}

func main() {

	p1 := &page{title: "test", body: []byte("this is a test page")}
	p1.save()

	//p2, err := loadPage(p1.title)
	p2, err := p1.load()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(p2.body))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit1/", editHandler1)
	http.HandleFunc("/edit2/", editHandler2)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
