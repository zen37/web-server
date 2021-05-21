package main

import "io/ioutil"

//describes how the wiki page will be stored in memory
type page struct {
	title string
	body  []byte // byte slice rather than string as this is the type expected by the io libraries
}

/*
save method that takes as its receiver a pointer to page
it returns an error value because that is the return type of WriteFile
if there is no error it returns nil, the zero-value for pointers, interfaces, and some other types
octal integer literal 0600, indicates that the file should be created with read-write permissions for the current user only.
*/
func (p *page) save() error {
	f := p.title + ".txt"
	return ioutil.WriteFile(f, p.body, 0600)
}

func main() {

}
