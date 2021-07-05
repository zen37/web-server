
package main

import (
	"fmt"
	"net/http"
)

func main() {
	h := Handler{
		Name: "Meowty Cat",
	}
	// http.Handle accepts an interface with the method ServeHTTP(w http.ResponseWriter, r *http.Request)
	http.Handle("/cat", h)

	// http.HandleFunc accepts an http.HandleFunc - a function that looks like Name(w http.ResponseWriter, r *http.Request)
	http.HandleFunc("/dog", HandlerFunc)

	http.ListenAndServe(":3000", nil)
}

type Handler struct {
	Name string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, this is ", h.Name)
}

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, this is Wooftastic Dog!")
}
When we call http.HandleFunc, we need to provide a function that takes two arguments: an http.ResponseWriter, and an *http.Request.

When we call http.Handle, we need to provide any type with the ServeHTTP(...) method. That ServeHTTP method is identical to the method passed into http.HandleFunc, which means the code above could actually ONLY use http.HandleFunc if it wanted to:

// We changed the second argument from "h" to "h.ServeHTTP"
http.HandleFunc("/cat", h.ServeHTTP)

http.HandleFunc("/dog", HandlerFunc)
After making this change our code will work exactly like the previous version. 🤯

If we try to do the opposite - to only use http.Handle - we will probably run into an error the first time. For instance, the following code results in an error.

func main() {
	h := Handler{
		Name: "Meowty Cat",
	}
	http.Handle("/cat", h)

	// This errors
	http.Handle("/dog", HandlerFunc)

	http.ListenAndServe(":3000", nil)
}
The exact error is below, or run it on the Go playground at https://play.golang.org/p/m2KEo8kRyWW​

cannot use HandlerFunc (type func(http.ResponseWriter, *http.Request)) as type http.Handler in argument to http.Handle:
func(http.ResponseWriter, *http.Request) does not implement http.Handler (missing ServeHTTP method)
From that you might conclude that we should always use http.HandleFunc - it works in both cases and http.Handle doesn't! - but that isn't entirely true. In Go we can make function types.

type HandlerFunc func(http.ResponseWriter, *http.Request)
And we can add methods to those types.

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)
}
Note: This can be REALLY weird, even for experienced developers, so don't worry if it is a bit confusing. Just try to keep reading as I try to connect the dots shortly.

​http.HandlerFunc is an example of this, and it is where all this sample code comes from. Interestingly, we can even our a function that matches an http.HandlerFunc into that type.

Now that we have a new type, and that new type actually has a ServeHTTP(...) method, we can assign any handler function to a variable of this type.

var hf http.HandlerFunc
hf = Dog

// Dog is an http handler function
func Dog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, this is Wooftastic Dog!")
}
And how the "hf" variable HAS a ServeHTTP method that just calls the original Dog function! 🤯

What this ultimately means is that we can ALSO write the original code using only http.Handle!

package main

import (
	"fmt"
	"net/http"
)

func main() {
	h := Handler{
		Name: "Meowty Cat",
	}
	http.Handle("/cat", h)

	var hf http.HandlerFunc
	hf = Dog

	http.Handle("/dog", hf)

	http.ListenAndServe(":3000", nil)
}

type Handler struct {
	Name string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, this is ", h.Name)
}

func Dog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, this is Wooftastic Dog!")
}
We can even take a shortcut and just convert the Dog function to the http.HandlerFunc type:

// No need to use the "var hf http.HandlerFunc" variable
http.Handle("/dog", http.HandlerFunc(Dog))
And now things should magically be clear to you ... just kidding! 😂

I get it. This is all confusing. Even though I spent all this time writing up an explanation, I still fully expect many readers to not quite get it. There are a lot of confusing techniques being used here that just won't click immediately for everyone. And that is the entire point of this particular course update! It IS confusing, and the course should spend more time trying to explain what is going on, showing errors, and helping students prepare for when they inevitably encounter a similar bug in their code.