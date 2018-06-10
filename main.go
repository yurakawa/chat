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
	"github.com/stretchr/objx"
	"github.com/yurakawa/trace"
)

// set the active Avatar implementation
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	err := t.templ.Execute(w, data)
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
		facebook.New("CLIENTID", "SECRET", "http://localhost:8080/auth/callback/facebook"),
		github.New("CLIENTID", "SECRET", "http://localhost:8080/auth/callback/github"),
		google.New("CLIENTID", "SECRET", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	// Open chat room
	go r.run()
	// Start web server
	log.Println("Start web Server. Port:", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
