package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/onufert/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handle the http request
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
	t.templ.Execute(w, data)
}

func main() {

	var addr = flag.String("addr", ":8080", "Address to serve application from")
	var debug = flag.Bool("debug", true, "Turn on debugging")
	flag.Parse()

	gomniauth.SetSecurityKey("some long key")
	gomniauth.WithProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	if *debug {
		r.tracer = trace.New(os.Stdout)
	}
	http.Handle("/static/",
		http.StripPrefix("/static",
			http.FileServer(http.Dir("static/"))))
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	// start the room for clients to connect to
	go r.run()
	//start the web server
	log.Println("Starting web server on port ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
