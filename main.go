package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/yurakawa/trace"
)

// templ represents a single templat
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServerHTTP is a http.HandleFunc that renders this template.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.templ.Execute(w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "Application Address")
	flag.Parse()

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("SECURITYKEY")
	gomniauth.WithProviders(
		facebook.New("CLIENTID", "SECRET", "http://localhost:8080/autn/callback/facebook"),
		github.New("CLIENTID", "SECRET", "http://localhost:8080/autn/callback/github"),
		google.New("CLIENTID", "SECRET", "http://localhost:8080/autn/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// Open chat room
	go r.run()
	// Start web server
	log.Println("Start web Server. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
