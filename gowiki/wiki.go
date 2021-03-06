package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

const (
	c_edit string = "edit.html"
	c_view string = "view.html"
)

//describes how the wiki page will be stored in memory
type page struct {
	Title string
	Body  []byte // byte slice rather than string as this is the type expected by the io libraries
}

/*
The function template.Must is a convenience wrapper that panics when passed a non-nil error value,
and otherwise returns the *Template unaltered. A panic is appropriate here;
if the templates can't be loaded the only sensible thing to do is exit the program.

The ParseFiles function takes any number of string arguments that identify our template files,
and parses those files into templates that are named after the base file name.
*/
var templates = template.Must(template.ParseFiles(c_edit, c_view))

/*
The function regexp.MustCompile will parse and compile the regular expression, and return a regexp.Regexp.
MustCompile is distinct from Compile in that it will panic if the expression compilation fails,
while Compile returns an error as a second parameter.
*/
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getPageTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	//fmt.Println("r.URL.Path = ", r.URL.Path)

	m := validPath.FindStringSubmatch(r.URL.Path)
	//fmt.Println("m = ", m)
	if m == nil {
		res := strings.Split(r.URL.Path, "/")
		t := res[len(res)-1]
		w.Write([]byte(t + " is an invalid Page Title"))
		//	http.NotFound(w, r)
		return "", errors.New(t + " is an invalid Page Title")
	}

	return m[2], nil //title is in second subexpression
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
	//t := r.URL.Path[len("/view/"):]
	t, err := getPageTitle(w, r)
	if err != nil {
		fmt.Println(err)
		//	w.Write([]byte(err.Error()))
		return
	}
	p, err := loadPage(t)
	if err != nil {
		// if the page doesn't exist the client is redirected to the edit Page so the content may be created
		http.Redirect(w, r, "/edit/"+t, http.StatusFound)
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

//method to load the page
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

	//t := r.URL.Path[len("/edit/"):]

	t, err := getPageTitle(w, r)
	if err != nil {
		fmt.Println(err)
		//	w.Write([]byte(err.Error()))
		return
	}

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

/*
There is an inefficiency in this code: renderTemplate calls ParseFiles every time a page is rendered.nb
A better approach would be to call ParseFiles once at program initialization,
parsing all templates into a single *Template.
Then we can use the ExecuteTemplate method to render a specific template.
*/
func renderTemplate(w http.ResponseWriter, p *page, tmpl string) {

	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	/*
		t, err := template.ParseFiles(tmpl + ".html")
		if err != nil {
			log.Fatal(err)
		}
		err = t.Execute(w, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/
}

//will handle the submission of the form located in the edit page
func saveHandler(w http.ResponseWriter, r *http.Request) {

	t, err := getPageTitle(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	b := r.FormValue("body") //returns string, that needs to be converted to []byte for page struct
	p := &page{Title: t, Body: []byte(b)}
	//saves data to a file
	err = p.save()
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
	http.HandleFunc("/edit/", editHandler2)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("Listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
