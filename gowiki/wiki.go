package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

//describes how the wiki page will be stored in memory
type page struct {
	Title string
	Body  []byte // byte slice rather than string as this is the type expected by the io libraries
}

/*
save method that takes as its receiver a pointer to page
it returns an error value because that is the return type of WriteFile
if there is no error it returns nil (the zero-value for pointers, interfaces, and some other types)
octal integer literal 0600, indicates that the file should be created with read-write permissions for the current user only.
*/
func (p *page) save() error {
	f := p.Title + ".txt"
	return ioutil.WriteFile(f, p.Body, 0600)
}

// why not a method? tutorial has a function - May 22
func loadPage(Title string) (*page, error) {
	f := Title + ".txt"
	Body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{Title: Title, Body: Body}, nil
}

// load method not in the tutorial - May 22
func (p *page) load() (*page, error) {
	f := p.Title + ".txt"
	Body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{Title: p.Title, Body: Body}, nil

}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	/*
		var p1 page
		p1.Title = r.URL.Path[len("/view/"):]
		p, err := p1.load()
	*/

	t := r.URL.Path[len("/view/"):]
	p, err := loadPage(t)
	if err != nil {
		// if the page doesn't exist the client is redirected to the edit Page so the content may be created
		http.Redirect(w, r, "/edit2/"+t, http.StatusFound)
		//The http.Redirect function adds an HTTP status code of http.StatusFound (302)
		//and a Location header to the HTTP response.
		return
	}
	// fmt.Fprintln(w, string(p2.Body))
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	// templ, err := template.ParseFiles("view.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// templ.Execute(w, p)

	renderTemplate(w, p, "view")
}

//The function editHandler* loads the page
//(or, if it doesn't exist, create an empty Page struct), and displays an HTML form.

//uses a method to load the page
func editHandler1(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Path[len("/edit1/"):]
	p1 := &page{Title: t}
	p2, err := p1.load()

	if err != nil {
		fmt.Fprintf(w, "<h1>editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"Body\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form>",
			p1.Title, p1.Title, p1.Body)
	} else {
		fmt.Fprintf(w, "<h1>editing %s</h1>"+
			"<form action=\"/save/%s\" method=\"POST\">"+
			"<textarea name=\"Body\">%s</textarea><br>"+
			"<input type=\"submit\" value=\"Save\">"+
			"</form>",
			p2.Title, p2.Title, p2.Body)
	}
	fmt.Fprintf(w, "<h6>func (p *page) load() (*page, error) has been used</h6>")
}

//uses a function to load the page
func editHandler2(w http.ResponseWriter, r *http.Request) {

	t := r.URL.Path[len("/edit2/"):]
	p, err := loadPage(t)
	if err != nil {
		p = &page{Title: t}
	}

	// will read the contents of edit.html and return a *template.Template.
	// templ, _ := template.ParseFiles("edit.html")
	// executes the template, writing the generated HTML to the http.ResponseWriter
	// templ.Execute(w, p)

	//fmt.Fprintf(w, "<h6>func loadPage(Title string) (*page, error) has been used</h6>")
	renderTemplate(w, p, "edit")
}

func renderTemplate(w http.ResponseWriter, p *page, tmpl string) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, p)
}

//will handle the submission of the form located in the edit page
func saveHandler(w http.ResponseWriter, r *http.Request) {

	//the title, provided in the URL, and teh form's obly field are stored in a new page
	t := r.URL.Path[len("/save/"):]
	b := r.FormValue("body") //returns string, that needs to be converted to []byte for page struct
	p := &page{Title: t, Body: []byte(b)}
	//saves data to a file
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//client is redirected
	http.Redirect(w, r, "/view/"+t, http.StatusFound)

}

func main() {

	p1 := &page{Title: "test", Body: []byte("this is a test page")}
	p1.save()

	//p2, err := loadPage(p1.Title)
	p2, err := p1.load()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit1/", editHandler1)
	http.HandleFunc("/edit2/", editHandler2)
	http.HandleFunc("/save/", saveHandler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
