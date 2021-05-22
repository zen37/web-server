package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

func loadPage(title string) (*page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{title: title, body: body}, nil
}

func (p *page) load() (*page, error) {
	f := p.title + ".txt"
	body, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return &page{title: p.title, body: body}, nil
}

func main() {
	p1 := &page{title: "test", body: []byte("this is a test page, loaded using a method")}
	p1.save()

	//p2, err := loadPage(p1.title)
	p2, err := p1.load()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(p2.body))
}
