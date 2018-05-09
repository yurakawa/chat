package main

import (
	"net/http"
	"log"
	"sync"
	"text/template"
	"path/filepath"
)


// templ represents a single templat
type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

// ServerHTTP is a http.HandleFunc that renders this template.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.templ.Execute(w, nil)
	if err  != nil {
		log.Fatal(err)
	}
}

func main() {
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// Start web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
